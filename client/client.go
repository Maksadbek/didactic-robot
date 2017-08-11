package client

import (
	"net"
	"net/http"
	"time"
)

type Config struct {
	ResolveTimeout      int
	TLSHandshakeTimeout int
}

type Client struct {
	httpClient *http.Client
}

type Dialer struct {
	*net.Dialer
}

func NewClient(config Config) (*Client, error) {
	dialer := &Dialer{
		Dialer: &net.Dialer{
			Timeout: time.Millisecond * time.Duration(config.ResolveTimeout),
		},
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext:         dialer.DialContext,
			TLSHandshakeTimeout: time.Millisecond * time.Duration(config.TLSHandshakeTimeout),
		},
	}

	return &Client{httpClient: &httpClient}, nil
}
