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

func init() {
	sig = make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
}

func SearchGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := validateQueryValues(r.URL.Query()); err != nil {
		statusBadRequest(w, r)
		return
	}

	searchByText(w, r)
}

var sig chan os.Signal

const (
	maxRemoteConns = 5

	socketReadSize  = 160
	socketWriteSize = 800

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
	connsRemote := connections.RemoteConnections[remoteAddr]
	connections.RUnlock()
	if connsRemote >= maxRemoteConns {
		http.Error(w, "Limit of simultaneous connections per remote address exceeded", http.StatusTooManyRequests)
		logger.Errorf("Simultaneous connections exceeded. Remote: %s", remoteAddr)
		return
	}

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
		select {
		case <-socket.Close:
			return
		case <-sig:
			msg := websocket.FormatCloseMessage(websocket.CloseServiceRestart, "Server shutdown")
			if err := socket.Conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
				logger.Errorf("Close message could not be sent: %s", err)
			}
		}
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
	ID        int64        `json:"id"`
	Term      string       `json:"term"`
	Filter    socketFilter `json:"filter"`
	Items     bool         `json:"items"`
	Locations bool         `json:"locations"`
	Features  bool         `json:"features"`
}

type socketFilter struct {
	Category string `json:"item,omitempty"`
	Location string `json:"location,omitempty"`
}

type socketResponse struct {
	ID    int64                 `json:"id"`
	Items []*model.SearchResult `json:"items,omitempty"`
	Error interface{}           `json:"error"`
}

func (s *socket) read(remote string) {
	defer func() {
		close(s.Close)
		close(s.Send)
		if err := s.Conn.Close(); err != nil {
			logger.Errorf("Error while closing socket: %s", err)
		}

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
		if err := s.Conn.SetReadDeadline(time.Now().Add(socketReadDeadline)); err != nil {
			logger.Errorf("Deadline could not be set: %s", err)
		}

		if atomic.LoadUint64(&requests) >= s.MaxRequests {
			msg := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Rate limit exceeded")
			if err := s.Conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
				logger.Errorf("Close message could not be sent: %s", err)
			}
			break
		}

		atomic.AddUint64(&requests, 1)

		req := &socketRequest{}
		if err := s.Conn.ReadJSON(req); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseServiceRestart) {
				logger.Error(err)
				break
			}

			if strings.HasPrefix(err.Error(), "invalid character") {
				logger.Errorf("Invalid data over socket: %v", err)

				msg := websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "Invalid data")
				if err := s.Conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
					logger.Errorf("Close message could not be sent: %s", err)
				}
			}

			break
		}

		go func() {
			res := &socketResponse{
				ID:    req.ID,
				Items: make([]*model.SearchResult, 0, 10),
			}

			q := req.Term

			if err := validateTerm(q); err != nil {
				res.Error = err.Error()
				s.Send <- res
				return
			}

			filter := &model.SearchFilter{
				Category: req.Filter.Category,
				Location: req.Filter.Location,
				ByName:   true,
			}

			search := model.NewSearch(cleanupString(q), filter, 5)

			if req.Items {
				search.Tasks.Add(1)
				go search.Items()
			}
			if req.Locations {
				search.Tasks.Add(1)
				go search.Locations()
			}
			if req.Features {
				search.Tasks.Add(1)
				go search.Features()
			}

			search.Close()

			for r := range search.Results {
				res.Items = append(res.Items, r...)
			}

			if err := search.Error; err != nil {
				logger.Error(err)

				var msg []byte
				switch {
				case errors.Is(err, item.ErrInvalidCategory):
					res.Error = err.Error()
					s.Send <- res
					return
				case errors.Is(err, api.ErrUnreachable):
					msg = websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "Service currently unavailable")
				default:
					msg = websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Internal server error")
				}

				if err := s.Conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
					logger.Errorf("Close message could not be sent: %s", err)
				}

				return
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

			if err := s.Conn.SetWriteDeadline(time.Now().Add(socketWriteDeadline)); err != nil {
				logger.Errorf("Deadline could not be set: %s", err)
			}

			if err := s.Conn.WriteJSON(res); err != nil {
				logger.Error(err)
			}
		case <-s.Close:
			return
		}
	}
}

func searchByText(w http.ResponseWriter, r *http.Request) {
	query := cleanupString(r.FormValue("query"))

	f, t, err := getFilter(query)
	if err != nil {
		statusBadRequest(w, r)
		return
	}

	if err := validateTerm(t); err != nil {
		statusBadRequest(w, r)
		return
	}

	p, err := model.CreatePageWithAPI(r.URL)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	search := model.NewSearch(t, f, 50)

	search.Tasks.Add(2)
	go search.Items()
	go search.Locations()

	search.Close()

	result := make([]*model.SearchResult, 0, 50)
	for r := range search.Results {
		result = append(result, r...)
	}

	if err := search.Error; err != nil {
		getErrorStatus(err, w, r)
		return
	}

	data, err := p.Result(result, query)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	view.RenderHTML("list", data, w)
}

func getFilter(t string) (*model.SearchFilter, string, error) {
	filter := &model.SearchFilter{}

	if kv := strings.Split(t, ":"); len(kv) == 2 {
		v := strings.TrimSpace(kv[1])
		switch k := kv[0]; k {
		case "item":
			if vt := strings.SplitN(v, " ", 2); len(vt) >= 2 {
				filter.Category = strings.TrimSpace(vt[0])
				t = strings.TrimSpace(vt[1])
			} else {
				return filter, "", errors.New("term is missing")
			}
		default:
			return filter, "", fmt.Errorf("unknown filter key \"%s\"", k)
		}
	}

	return filter, t, nil
}
