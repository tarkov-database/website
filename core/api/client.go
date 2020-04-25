package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const contentTypeJSON = "application/json"

var (
	userAgent = "tarkov-database-frontend"
	client    *http.Client
)

func request(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	u := cfg.URL
	if len(path) > 1 {
		u += path
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Token))

	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		switch {
		case strings.HasSuffix(err.Error(), context.DeadlineExceeded.Error()), strings.HasSuffix(err.Error(), "connection refused"):
			return res, ErrUnreachable
		default:
			return res, err
		}
	}

	if res.Header.Get("Content-Type") != contentTypeJSON {
		return res, ErrWrongContentType
	}

	return res, nil
}

func decodeBody(body io.ReadCloser, target interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(target)
}

func encodeBody(w io.Writer, source interface{}) error {
	return json.NewEncoder(w).Encode(source)
}

type Options struct {
	Sort   string
	Limit  int
	Offset int
	Filter map[string]string
}

func GET(ctx context.Context, path string, opts *Options, target interface{}) error {
	v := url.Values{}
	for key, val := range opts.Filter {
		if val != "" {
			v.Add(key, url.QueryEscape(val))
		}
	}
	if len(opts.Sort) > 0 {
		v.Add("sort", opts.Sort)
	}
	if opts.Limit > 0 {
		v.Add("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		v.Add("offset", strconv.Itoa(opts.Offset))
	}

	res, err := request(ctx, "GET", fmt.Sprintf("%s?%s", path, v.Encode()), nil)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return statusToError(res)
	}

	if err = decodeBody(res.Body, target); err != nil {
		return ErrParsing
	}

	return nil
}

func POST(ctx context.Context, path string, source interface{}) error {
	buf := new(bytes.Buffer)

	if err := encodeBody(buf, source); err != nil {
		return ErrParsing
	}

	res, err := request(ctx, "POST", path, buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return statusToError(res)
	}

	return nil
}

func PUT(ctx context.Context, path string, source interface{}) error {
	buf := new(bytes.Buffer)

	if err := encodeBody(buf, source); err != nil {
		return ErrParsing
	}

	res, err := request(ctx, "PUT", path, buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return statusToError(res)
	}

	return nil
}
