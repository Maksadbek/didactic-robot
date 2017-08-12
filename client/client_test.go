package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	config := Config{
		ResolveTimeout:      time.Second * time.Duration(3),
		TLSHandshakeTimeout: time.Second * 3,
		NameserverAddr:      "8.8.8.8:53",
		ConnectionTimeout:   time.Second * 3,
		TCPKeepIdleTime:     time.Second * 3,
		TCPKeepIntervalTime: time.Second * 3,
		TCPKeepFailAfter:    3,
		ReadTimeout:         time.Second * 60,
		WriteTimeout:        time.Second * 60,
	}

	client, err := NewClient(config)
	require.NoError(t, err)

	t.Run("IP address", func(t *testing.T) {
		testMsg := "hello"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, testMsg)
		}))

		defer server.Close()

		request, err := http.NewRequest(http.MethodGet, server.URL, nil)
		require.NoError(t, err)

		resp, err := client.Do(request)
		require.NoError(t, err)

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, string(body), testMsg)
	})

	t.Run("host", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "http://httpbin.org/get", nil)
		require.NoError(t, err)

		resp, err := client.Do(request)
		require.NoError(t, err)

		defer resp.Body.Close()

		_, err = ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
	})

	t.Run("connection timeout", func(t *testing.T) {
		var err error
		var request *http.Request

		time.AfterFunc(config.ConnectionTimeout+time.Second*10, func() {
			require.NotNil(t, err)
			require.Equal(t, err.Error(), "Get http://111.111.111.111: dial tcp 111.111.111.111:80: i/o timeout")
		})

		request, err = http.NewRequest(http.MethodGet, "http://111.111.111.111" /* unreachable address */, nil)
		require.NoError(t, err)

		_, err = client.Do(request)
	})

	t.Run("https", func(t *testing.T) {
		var err error
		var request *http.Request

		var buf = &bytes.Buffer{}

		blob := "{'foo':'bar'}"
		buf.WriteString(blob)

		request, err = http.NewRequest(http.MethodPost, "https://httpbin.org/post", buf)
		require.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(request)
		require.NoError(t, err)

		defer resp.Body.Close()

		m := struct {
			Data string `json:"data"`
		}{}

		err = json.NewDecoder(resp.Body).Decode(&m)
		require.NoError(t, err)

		require.Equal(t, m.Data, blob)
	})
}
