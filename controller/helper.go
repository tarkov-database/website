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

	"github.com/google/logger"
)

var (
	ErrTooLongShort = errors.New("string is too short or too long")
	ErrIllegalChars = errors.New("string has illegal characters")
)

var hostname string

func init() {
	if env := os.Getenv("HOST"); len(env) > 0 {
		hostname = env
	} else {
		logger.Warning("Host is not set!")
	}
}

func getPage(u *url.URL) int {
	page := u.Query().Get("p")
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

var regexNonAlnum = regexp.MustCompile(`[^[:alnum:]]`)

func isAlnum(s string) bool {
	return !regexNonAlnum.MatchString(s)
}

var regexNonAlnumExtended = regexp.MustCompile(`[^[:alnum:][:blank:]!#$%&'()*+,\-./:;?_~]`)

func isAlnumExtended(s string) bool {
	return !regexNonAlnumExtended.MatchString(s)
}

func validateQueryValues(q url.Values) error {
	for k, v := range q {
		if len(k) > 30 {
			return fmt.Errorf("error in key \"%s\": %w", k, ErrTooLongShort)
		}
		if !isAlnum(k) {
			return fmt.Errorf("error in key \"%s\": %w", k, ErrIllegalChars)
		}
		for _, e := range v {
			v, err := url.QueryUnescape(e)
			if err != nil {
				return fmt.Errorf("error in value of \"%s\": %w", k, err)
			}
			if len(v) > 100 {
				return fmt.Errorf("error in value of \"%s\": %w", k, ErrTooLongShort)
			}
			if !isAlnumExtended(v) {
				return fmt.Errorf("error in value of \"%s\": %w", k, ErrIllegalChars)
			}
		}
	}

	return nil
}

func unescapeParams(params map[string]string) error {
	for k, v := range params {
		var err error
		params[k], err = url.QueryUnescape(v)
		if err != nil {
			return fmt.Errorf("error in value of \"%s\": %w", k, err)
		}
	}

	return nil
}

func validateTerm(q string) error {
	if len(q) < 3 || len(q) > 32 {
		return ErrTooLongShort
	}
	if !isAlnumExtended(q) {
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

func isSupportedMediaType(r *http.Request) bool {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		return true
	case "application/geo+json":
		return true
	}

	return false
}
