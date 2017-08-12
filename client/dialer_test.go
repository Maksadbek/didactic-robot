package client

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDialer(t *testing.T) {
	dialer := &Dialer{
		Dialer:           &net.Dialer{},
		nameserverAddr:   "8.8.8.8:53",
		keepIdleTime:     3,
		keepIntervalTime: 3,
		keepFailAfter:    3,
		readTimeout:      60,
		writeTimeout:     60,
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
		},
	}

	resp, err := client.Get("http://httpbin.org/get")
	require.NoError(t, err)

	defer resp.Body.Close()

	require.NotNil(t, resp.Body)
}
