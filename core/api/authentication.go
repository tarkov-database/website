package api

import (
	"context"
	"log"
	"time"
)

func refreshScheduler() {
	if err := refreshToken(); err != nil {
		log.Printf("Error while refreshing token: %s", err)
		time.Sleep(3 * time.Second)
		go refreshScheduler()
		return
	}

	claims, err := cfg.GetTokenClaims()
	if err != nil {
		log.Fatalf("Error while getting token claims: %s", err)
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

	resp := &tokenResponse{}

	if err := GET(ctx, "/token", &Options{}, resp); err != nil {
		return err
	}

	cfg.Token = resp.Token

	return nil
}
