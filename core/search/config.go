package search

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/net/http2"
)

var config *clientConfig

func init() {
	var err error

	config, err = newConfig()
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(2)
	}

	if config.TLS {
		clientCert, rootCAs, err := config.GetTLSCertificate()
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(2)
		}

		tlsConfig := &tls.Config{RootCAs: rootCAs}
		if len(config.Certificate) > 0 {
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

	go refreshScheduler()
}

type clientConfig struct {
	Host        string
	Token       string
	TLS         bool
	Certificate string
	PrivateKey  string
	RootCAs     string
}

func newConfig() (*clientConfig, error) {
	c := &clientConfig{}

	if env := os.Getenv("SEARCH_HOST"); len(env) > 4 {
		c.Host = env
	} else {
		return c, errors.New("search host not set")
	}

	if env := os.Getenv("SEARCH_TOKEN"); len(env) > 4 {
		c.Token = env
	} else {
		return c, errors.New("search token not set")
	}

	if env := os.Getenv("SEARCH_TLS"); len(env) > 3 {
		if b, err := strconv.ParseBool(env); err == nil {
			c.TLS = b
		}

		if env := os.Getenv("SEARCH_CERT"); len(env) > 0 {
			c.Certificate = env

			if env := os.Getenv("SEARCH_KEY"); len(env) > 0 {
				c.PrivateKey = env
			} else {
				return c, errors.New("search private key not set")
			}
		}

		if env := os.Getenv("SEARCH_CA"); len(env) > 0 {
			c.RootCAs = env
		}
	}

	return c, nil
}

func (c *clientConfig) GetTLSCertificate() (tls.Certificate, *x509.CertPool, error) {
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
