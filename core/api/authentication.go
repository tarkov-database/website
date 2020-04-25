package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/logger"
)

func refreshScheduler() {
	if err := refreshToken(); err != nil {
		logger.Error(err)
		time.Sleep(3 * time.Second)
		go refreshScheduler()
		return
	}

	claims, err := cfg.GetTokenClaims()
	if err != nil {
		logger.Fatal(err)
		return
	}

	refresh := claims.ExpirationTime.Add(-30 * time.Second).Sub(time.Now())
	time.Sleep(refresh)
	go refreshScheduler()
}

type tokenResponse struct {
	Token string `json:"token"`
}

func refreshToken() error {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	res, err := request(ctx, "GET", "/token", nil)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return statusToError(res)
	}

	resp := tokenResponse{}

	if err = decodeBody(res.Body, &resp); err != nil {
		logger.Errorf("Error while parsing json: %s", err)
		return ErrParsing
	}

	if res.StatusCode == http.StatusCreated {
		cfg.Token = resp.Token
	}

	return nil
}
