package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrUnreachable      = errors.New("api: server unreachable")
	ErrAuthentication   = errors.New("api: authentication failed")
	ErrWrongContentType = errors.New("api: wrong content type")
	ErrParsing          = errors.New("api: parsing error")
)

type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
}

func (e Response) Error() string {
	return fmt.Sprintf("%v: %v", e.StatusCode, e.Message)
}

func getStatus(res *http.Response) (*Response, error) {
	e := &Response{}
	if err := decodeBody(res.Body, e); err != nil {
		return e, ErrParsing
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

type Timestamp struct {
	time.Time
}

func (u *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Unix())
}

func (u *Timestamp) UnmarshalJSON(b []byte) error {
	var i int64

	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	*u = Timestamp{time.Unix(i, 0)}

	return nil
}
