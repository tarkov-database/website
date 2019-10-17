package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/logger"
)

var tokenExp int64

func setTokenExpiration() {
	claims, err := cfg.GetTokenClaims()
	if err != nil {
		logger.Fatal(err)
		return
	}
	tokenExp = int64(claims["exp"].(float64)) - 20
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

	var j map[string]interface{}
	if err = decodeBody(res.Body, &j); err != nil {
		logger.Errorf("Error while parsing json: %s", err)
		return ErrParsing
	}

	if token, ok := j["token"].(string); res.StatusCode == http.StatusCreated && ok {
		cfg.Token = token
	}

	setTokenExpiration()

	return nil
}
