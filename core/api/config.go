package api

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/http2"
)

var cfg *config

func init() {
	var err error

	cfg, err = newConfig()
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(2)
	}

	if cfg.TLS {
		clientCert, rootCAs, err := cfg.GetTLSCertificate()
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(2)
		}

		tlsConfig := &tls.Config{RootCAs: rootCAs}
		if len(cfg.Certificate) > 0 {
			tlsConfig.Certificates = []tls.Certificate{clientCert}
		}

		client = &http.Client{Transport: &http2.Transport{
			TLSClientConfig: tlsConfig,
		}}
	} else {
		client = &http.Client{Transport: &http.Transport{
			IdleConnTimeout: 0,
		}}
	}

	userAgent = fmt.Sprintf("%s", userAgent)

	go refreshScheduler()
}

type config struct {
	URL         string
	Token       string
	TLS         bool
	Certificate string
	PrivateKey  string
	RootCAs     string
}

func newConfig() (*config, error) {
	c := &config{}

	if env := os.Getenv("API_URL"); len(env) > 4 {
		c.URL = env
	} else {
		return c, errors.New("api url not set")
	}

	if env := os.Getenv("API_TOKEN"); len(env) > 4 {
		c.Token = env
	} else {
		return c, errors.New("api token not set")
	}

	if env := os.Getenv("API_TLS"); len(env) > 3 {
		if b, err := strconv.ParseBool(env); err == nil {
			c.TLS = b
		}

		if env := os.Getenv("API_CERT"); len(env) > 0 {
			c.Certificate = env

			if env := os.Getenv("API_KEY"); len(env) > 0 {
				c.PrivateKey = env
			} else {
				return c, errors.New("api private key not set")
			}
		}

		if env := os.Getenv("API_CA"); len(env) > 0 {
			c.RootCAs = env
		}
	}

	return c, nil
}

func (c *config) GetTLSCertificate() (tls.Certificate, *x509.CertPool, error) {
	var clientCert tls.Certificate
	var rootCAs *x509.CertPool

	if len(c.Certificate) > 0 {
		clientCertPEM, err := ioutil.ReadFile(c.Certificate)
		if err != nil {
			return clientCert, rootCAs, err
		}

		clientKeyPEM, err := ioutil.ReadFile(c.PrivateKey)
		if err != nil {
			return clientCert, rootCAs, err
		}

		clientCert, err = tls.X509KeyPair(clientCertPEM, clientKeyPEM)
		if err != nil {
			return clientCert, rootCAs, err
		}
	}

	if len(c.RootCAs) > 0 {
		rootCAs = x509.NewCertPool()

		caPEM, err := ioutil.ReadFile(c.RootCAs)
		if err != nil {
			return clientCert, rootCAs, err
		}

		if ok := rootCAs.AppendCertsFromPEM(caPEM); !ok {
			return clientCert, rootCAs, errors.New("failed to load root CA")
		}
	}

	return clientCert, rootCAs, nil
}

type tokenClaims struct {
	ExpirationTime Timestamp `json:"exp"`
	NotBefore      Timestamp `json:"nbf"`
	IssuedAt       Timestamp `json:"iat"`
}

func (c *config) GetTokenClaims() (*tokenClaims, error) {
	claims := &tokenClaims{}

	seg := strings.Split(c.Token, ".")[1]
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	data, err := base64.URLEncoding.DecodeString(seg)
	if err != nil {
		return claims, err
	}

	if err = json.Unmarshal(data, &claims); err != nil {
		return claims, err
	}

	return claims, nil
}
