package controller

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tarkov-database/website/model/item"

	"github.com/google/logger"
)

var (
	ErrTooLongShort = errors.New("keyword is too short or too long")
	ErrIllegalChars = errors.New("keyword has illegal characters")
)

var hostname string

var itemKinds map[string]item.Kind

func init() {
	if env := os.Getenv("HOST"); len(env) > 0 {
		hostname = env
	} else {
		logger.Warning("Host is not set!")
	}
}

func getPage(r *http.Request) int {
	page := r.URL.Query().Get("p")
	var p int
	if len(page) != 0 {
		p, _ = strconv.Atoi(page)
	}

	return p
}

func getLimitOffset(page int) (limit, offset int) {
	offset, limit = 0, 100
	if page > 1 {
		page--
		offset = page * limit
	}
	return
}

func cleanupString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

var regexNonASCII = regexp.MustCompile(`[^[:ascii:]]`)

func isASCII(s string) bool {
	return !regexNonASCII.MatchString(s)
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

func getRemoteAddr(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastAddr := strings.TrimSpace(addrs[len(addrs)-1])
		ip := net.ParseIP(lastAddr)

		if ip == nil || ip.IsUnspecified() {
			logger.Errorf("Invalid XXF IP: %s", lastAddr)
		}

		return ip.String()
	}

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		ip := net.ParseIP(host)
		if ip == nil || ip.IsUnspecified() {
			logger.Errorf("Invalid remote IP: %s", host)
		}

		return ip.String()
	}

	return r.RemoteAddr
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if len(origin) == 0 {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	originHost := u.Hostname()
	if port := u.Port(); len(port) > 0 {
		originHost += ":" + u.Port()
	}

	return originHost == hostname
}

func getItemKind(id string) (item.Kind, error) {
	v, ok := itemKinds[id]
	if !ok {
		return v, errors.New("item not found")
	}

	return v, nil
}

func isSupportedMediaType(r *http.Request) bool {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		return true
	case "application/geo+json":
		return true
	}

	return false
}

type timingMetrics map[string]time.Duration

func addTimingHeader(metrics timingMetrics, w http.ResponseWriter) {
	v := make([]string, 0, len(metrics))
	for k, d := range metrics {
		v = append(v, fmt.Sprintf("%s;dur=%.3f", k, float64(d.Microseconds())/1000))
	}

	w.Header().Set("Server-Timing", strings.Join(v, ","))
}
