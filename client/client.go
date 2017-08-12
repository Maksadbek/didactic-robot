package client

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type Config struct {
	ResolveTimeout      time.Duration
	TLSHandshakeTimeout time.Duration
	ConnectionTimeout   time.Duration

	TCPKeepIdleTime     time.Duration
	TCPKeepIntervalTime time.Duration
	TCPKeepFailAfter    time.Duration

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	NameserverAddr string
}

type Client struct {
	httpClient *http.Client
	config     Config
}

func NewClient(config Config) (*Client, error) {
	dialer := &Dialer{
		keepIdleTime:     seconds(config.TCPKeepIdleTime),
		keepIntervalTime: seconds(config.TCPKeepIntervalTime),
		keepFailAfter:    seconds(config.TCPKeepFailAfter),
		readTimeout:      seconds(config.WriteTimeout),
		writeTimeout:     seconds(config.WriteTimeout),
		nameserverAddr:   config.NameserverAddr,

		Dialer: &net.Dialer{
			Timeout: config.ConnectionTimeout,
		},
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext:         dialer.DialContext,
			TLSHandshakeTimeout: time.Millisecond * time.Duration(config.TLSHandshakeTimeout),
		},
	}

	return &Client{
		httpClient: &httpClient,
		config:     config,
	}, nil
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.config.ResolveTimeout)
	defer cancel()

	return c.httpClient.Do(r)
}

func getRandomAddr(addrs []string) string {

	n := rand.Int31n(int32(len(addrs)))

	return addrs[n]
}

func init() {
	rand.Seed(time.Now().Unix())
}
