package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/logger"
)

var (
	ErrUnreachable    = errors.New("api: server unreachable")
	ErrAuthentication = errors.New("api: authentication failed")
	ErrParsing        = errors.New("api: parsing error")
)

type response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
}

func (e *response) Error() error {
	return fmt.Errorf("%v: %v", e.StatusCode, e.Message)
}

func getStatus(res *http.Response) (*response, error) {
	e := &response{}
	if err := decodeBody(res.Body, e); err != nil {
		logger.Errorf("Error while parsing json: %s", err)
		return e, ErrParsing
	}

	return e, nil
}

func statusToError(res *http.Response) error {
	r, err := getStatus(res)
	if err != nil {
		return err
	}

	return r.Error()
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
