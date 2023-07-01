package search

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

var (
	ErrInvalidTerm      = errors.New("given term is invalid")
	ErrUnreachable      = errors.New("server unreachable")
	ErrWrongContentType = errors.New("wrong content type")
	ErrParsing          = errors.New("parsing error")
)

const contentTypeJSON = "application/json"

var (
	userAgent = "tarkov-database-frontend"
	client    *http.Client
)

func request(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	u := config.Host
	if len(path) > 1 {
		u += path
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Token))

	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), strings.Contains(err.Error(), "connect:"):
			return res, ErrUnreachable
		default:
			return res, err
		}
	}

	if res.Header.Get("Content-Type") != contentTypeJSON {
		return res, ErrWrongContentType
	}

	if res.StatusCode >= 300 {
		return res, statusToError(res)
	}

	return res, nil
}

func decodeBody(body io.ReadCloser, target interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(target)
}

type StatusResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e StatusResponse) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Message)
}

func getStatus(res *http.Response) (*StatusResponse, error) {
	e := &StatusResponse{}
	if err := decodeBody(res.Body, e); err != nil {
		return e, fmt.Errorf("%w: %s", ErrParsing, err)
	}

	return e, nil
}

func statusToError(res *http.Response) error {
	r, err := getStatus(res)
	if err != nil {
		return err
	}

	return r
}

type Result struct {
	Count int64  `json:"count"`
	Data  []Item `json:"data"`
}

type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
}

type Query struct {
	Query string
}

type Options struct {
	Fuzzy bool
	Limit int
}

func Search(ctx context.Context, query *Query, opts *Options) (*Result, error) {
	result := &Result{}

	if len(query.Query) < 3 {
		return result, fmt.Errorf("%w: query is too short", ErrInvalidTerm)
	}

	v := url.Values{}

	v.Add("q", query.Query)

	if opts.Limit > 0 {
		v.Add("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Fuzzy {
		v.Add("fuzzy", strconv.FormatBool(opts.Fuzzy))
	}

	path := "/search"
	if len(v) > 0 {
		path = path + "?" + v.Encode()
	}

	res, err := request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, fmt.Errorf("GET \"%s\" %w", path, err)
	}

	if err = decodeBody(res.Body, result); err != nil {
		return result, fmt.Errorf("GET \"%s\" %w: %s", path, ErrParsing, err)
	}

	return result, nil
}
