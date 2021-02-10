package api

import (
	"context"
	"log"
	"time"
)

func refreshScheduler() {
	exp, err := refreshToken()

	if err != nil {
		log.Printf("Error while refreshing token: %s", err)
		time.Sleep(5 * time.Second)
		go refreshScheduler()
		return
	}

	refresh := exp.Add(-60 * time.Second).Sub(time.Now())

	time.Sleep(refresh)

	go refreshScheduler()
}

type tokenResponse struct {
	Token      string    `json:"token"`
	Expiration Timestamp `json:"expires"`
}

func refreshToken() (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	resp := &tokenResponse{}

	if err := GET(ctx, "/token", &Options{}, resp); err != nil {
		return time.Time{}, err
	}

	cfg.Token = resp.Token

	return resp.Expiration.Time, nil
}
