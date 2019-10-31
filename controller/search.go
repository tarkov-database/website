package controller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model"
	"github.com/tarkov-database/website/model/item"
	"github.com/tarkov-database/website/view"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var (
	ErrTooLongShort = errors.New("keyword is too short or too long")
	ErrIllegalChars = errors.New("keyword has illegal characters")
)

func SearchGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	getQuery(w, r)
}

const (
	maxConnsRemote = 5

	socketReadSize  = 128
	socketWriteSize = 768

	socketReadDeadline  = 3 * time.Minute
	socketWriteDeadline = 20 * time.Second
)

var connections = socketConnections{
	RemoteConnections: make(map[string]uint, 0),
}

var socketUpgrader = websocket.Upgrader{
	ReadBufferSize:  socketReadSize,
	WriteBufferSize: socketWriteSize,
	CheckOrigin:     checkOrigin,
}

func SearchWS(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	remoteAddr := getRemoteAddr(r)
	if remoteAddr == r.RemoteAddr {
		http.Error(w, "Remote address is not valid", http.StatusInternalServerError)
		logger.Errorf("Invalid remote address: %s", remoteAddr)
		return
	}

	connections.RLock()
	if connections.RemoteConnections[remoteAddr] >= maxConnsRemote {
		connections.RUnlock()
		http.Error(w, "Limit of simultaneous connections per remote address exceeded", http.StatusTooManyRequests)
		logger.Errorf("Simultaneous connections exceeded. Remote: %s", remoteAddr)
		return
	}
	connections.RUnlock()

	c, err := socketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("websocket upgrade: %s", err)
		return
	}

	c.SetReadLimit(socketReadSize)
	socket := &socket{
		MaxRequests: 10,
		Conn:        c,
		Send:        make(chan *socketResponse, 1),
		Close:       make(chan bool, 1),
	}

	connections.Lock()
	connections.RemoteConnections[remoteAddr]++
	connections.Unlock()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		msg := websocket.FormatCloseMessage(websocket.CloseServiceRestart, "Server shutdown")
		socket.Conn.WriteMessage(websocket.CloseMessage, msg)
	}()

	go socket.read(remoteAddr)
	go socket.write()
}

type socket struct {
	MaxRequests uint64
	Conn        *websocket.Conn
	Send        chan *socketResponse
	Close       chan bool
}

type socketConnections struct {
	sync.RWMutex
	RemoteConnections map[string]uint
}

type socketRequest struct {
	ID   int64  `json:"id"`
	Text string `json:"text"`
}

type socketResponse struct {
	ID    int64           `json:"id"`
	Items []*socketResult `json:"items"`
	Error interface{}     `json:"error"`
}

type socketResult struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Category  string `json:"category"`
}

func (s *socket) read(remote string) {
	defer func() {
		close(s.Close)
		close(s.Send)
		s.Conn.Close()

		time.Sleep(10 * time.Second)

		connections.Lock()
		connections.RemoteConnections[remote]--
		if connections.RemoteConnections[remote] == 0 {
			delete(connections.RemoteConnections, remote)
		}
		connections.Unlock()
	}()

	var requests uint64
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

	Loop:
		for {
			select {
			case <-s.Close:
				break Loop
			case <-ticker.C:
				val := atomic.LoadUint64(&requests)
				if val >= s.MaxRequests {
					atomic.AddUint64(&requests, ^uint64(s.MaxRequests-1))
				}
				if val < s.MaxRequests && val > 0 {
					atomic.AddUint64(&requests, ^uint64(val-1))
				}
			}
		}
	}()

	for {
		s.Conn.SetReadDeadline(time.Now().Add(socketReadDeadline))

		if atomic.LoadUint64(&requests) >= s.MaxRequests {
			msg := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Rate limit exceeded")
			s.Conn.WriteMessage(websocket.CloseMessage, msg)
			break
		}

		atomic.AddUint64(&requests, 1)

		req := &socketRequest{}
		err := s.Conn.ReadJSON(req)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logger.Error(err)
				break
			}
			if strings.HasPrefix(err.Error(), "invalid character") {
				logger.Errorf("Invalid data over socket: %s", err)
				msg := websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "Invalid data")
				s.Conn.WriteMessage(websocket.CloseMessage, msg)
			}
			break
		}

		go func() {
			res := &socketResponse{ID: req.ID}

			q := req.Text

			if err := validateKeyword(q); err != nil {
				res.Error = err
				s.Send <- res
				return
			}

			q = cleanupString(q)

			if op, _ := getOperator(q); op != "" {
				return
			}

			result, err := item.GetItemsBySearch(q, 5)
			if err != nil {
				logger.Error(err)

				var msg []byte
				switch err {
				case api.ErrUnreachable:
					msg = websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "Service currently unavailable")
				default:
					msg = websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Internal server error")
				}
				s.Conn.WriteMessage(websocket.CloseMessage, msg)

				return
			}

			items := result.GetEntities()

			res.Items = make([]*socketResult, len(items))
			for i, r := range items {
				cat, err := item.KindToCategory(r.GetKind())
				if err != nil {
					res.Error = err
					s.Send <- res
					return
				}

				res.Items[i] = &socketResult{
					ID:        r.GetID(),
					Name:      r.GetName(),
					ShortName: r.GetShortName(),
					Category:  strings.ReplaceAll(cat, " ", "-"),
				}
			}

			select {
			case <-s.Close:
				return
			default:
				s.Send <- res
			}
		}()
	}
}

func (s *socket) write() {
	for {
		select {
		case res, ok := <-s.Send:
			if !ok {
				return
			}

			s.Conn.SetWriteDeadline(time.Now().Add(socketWriteDeadline))

			if err := s.Conn.WriteJSON(res); err != nil {
				logger.Error(err)
			}
		case <-s.Close:
			return
		}
	}
}

func searchByText(kw string, w http.ResponseWriter, r *http.Request) {
	result, err := item.GetItemsBySearch(kw, 60)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	p, err := model.CreatePageWithAPI(r.URL)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	view.Render("list", p.ItemResult(result, kw, true), w)
}

func getOperator(q string) (operator string, query string) {
	if parts := strings.Split(q, ":"); len(parts) >= 2 {
		switch parts[0] {
		case "category", "cat":
			operator = "category"
			query = strings.ToLower(strings.TrimSpace(parts[1]))
		}
	}

	return
}

func getQuery(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("query")
	if err := validateKeyword(q); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	q = cleanupString(q)

	switch op, query := getOperator(q); op {
	case "category":
		http.Redirect(w, r, fmt.Sprintf("/item/%s", query), http.StatusMovedPermanently)
	default:
		searchByText(q, w, r)
	}
}

func validateKeyword(q string) error {
	if len(q) < 3 || len(q) > 32 {
		return ErrTooLongShort
	}
	if !isASCII(q) {
		return ErrIllegalChars
	}

	return nil
}
