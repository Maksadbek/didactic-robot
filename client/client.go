package client

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

type Config struct {
	ResolveTimeout      int
	TLSHandshakeTimeout int
	NameserverAddr      string
	ConnectionTimeout   time.Duration
	TCPKeepIdleTime     time.Duration
	TCPKeepIntervalTime time.Duration
	TCPKeepFailAfter    time.Duration

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
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

func (c *Client) Do(r *http.Request) (*http.Response, error) {

	ctx := context.Background()

	addrs, err := net.DefaultResolver.LookupHost(ctx, r.URL.Host)

	var addr string
	// if address is not found using local resolver
	// get address from remote name server
	if err != nil || len(addrs) == 0 {
		client := new(dns.Client)
		msg := new(dns.Msg)

		msg.SetQuestion(dns.Fqdn(r.URL.Host), dns.TypeA)

		reply, _, err := client.Exchange(msg, c.config.NameserverAddr)
		if err != nil {
			return nil, nil
		}

		if reply.Rcode != dns.RcodeSuccess {
			return nil, nil
		}

		for _, a := range reply.Answer {
			addr = dns.Field(a, 1)
			break
		}

		r.URL.Host = addr

		return c.httpClient.Do(r)
	}

	// if we have multiple addresses
	// choose one randomly
	if len(addrs) > 1 {
		r.URL.Host = getRandomAddr(addrs)
	} else {
		r.URL.Host = addrs[0]
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		r.URL.Host = getRandomAddr(addrs)
	} else {
		return resp, err
	}

	return c.httpClient.Do(r)
}

func getRandomAddr(addrs []string) string {
	rand.Seed(time.Now().Unix())

	n := rand.Int31n(int32(len(addrs)))

	return addrs[n]
}
