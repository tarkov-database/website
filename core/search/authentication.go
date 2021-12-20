package search

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tarkov-database/website/core/api"
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
	Token      string        `json:"token"`
	Expiration api.Timestamp `json:"expiresAt"`
}

func refreshToken() (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	tokenResp := &tokenResponse{}

	path := "/token"

	res, err := request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return time.Time{}, fmt.Errorf("GET \"%s\" %w", path, err)
	}

	if err = decodeBody(res.Body, tokenResp); err != nil {
		return time.Time{}, fmt.Errorf("GET \"%s\" %w: %s", path, ErrParsing, err)
	}

	config.Token = tokenResp.Token

	return tokenResp.Expiration.Time, nil
}
