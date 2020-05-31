package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/http2"
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

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), strings.Contains(err.Error(), "connect:"):
			return resp, ErrUnreachable
		default:
			return resp, err
		}
	}

	if resp.Header.Get("Content-Type") != contentTypeJSON {
		return resp, ErrWrongContentType
	}

	if resp.StatusCode >= 300 {
		return resp, statusToError(resp)
	}

	return resp, nil
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
	if opts.Sort != "" {
		v.Add("sort", opts.Sort)
	}
	if opts.Limit > 0 {
		v.Add("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		v.Add("offset", strconv.Itoa(opts.Offset))
	}

	if len(v) > 0 {
		path = path + "?" + v.Encode()
	}

	res, err := request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return fmt.Errorf("GET \"%s\" %w", path, err)
	}

	if err = decodeBody(res.Body, target); err != nil {
		// Retry after server closes the TCP connection after sending a GOAWAY frame with code 0.
		// Workaround for a probably faulty HTTP2 client or server implementation.
		// See https://github.com/golang/go/issues/18112, https://github.com/minio/minio/issues/7271
		var goAway http2.GoAwayError
		if errors.As(err, &goAway) && goAway.ErrCode == http2.ErrCodeNo {
			res, err := request(ctx, http.MethodGet, path, nil)
			if err != nil {
				return fmt.Errorf("GET \"%s\" %w", path, err)
			}

			if err = decodeBody(res.Body, target); err != nil {
				return fmt.Errorf("GET \"%s\" %w", path, err)
			}
		}

		return fmt.Errorf("GET \"%s\" %w: %s", path, ErrParsing, err)
	}

	return nil
}

func POST(ctx context.Context, path string, source interface{}) error {
	buf := new(bytes.Buffer)

	if err := encodeBody(buf, source); err != nil {
		return fmt.Errorf("POST \"%s\" %w: %s", path, ErrParsing, err)
	}

	res, err := request(ctx, http.MethodPost, path, buf)
	if err != nil {
		return fmt.Errorf("POST \"%s\" %w", path, err)
	}
	defer res.Body.Close()

	return nil
}

func PUT(ctx context.Context, path string, source interface{}) error {
	buf := new(bytes.Buffer)

	if err := encodeBody(buf, source); err != nil {
		return fmt.Errorf("PUT \"%s\" %w: %s", path, ErrParsing, err)
	}

	res, err := request(ctx, http.MethodPut, path, buf)
	if err != nil {
		return fmt.Errorf("PUT \"%s\" %w", path, err)
	}
	defer res.Body.Close()

	return nil
}
