package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/maksadbek/didactic-robot/client"
)

// flags
var (
	resolveTimeout      = flag.Duration("resolve-timeout", 0, "set resolve timeout")
	tlsHandshakeTimeout = flag.Duration("tls-handshake-timeout", 0, "set TLS handshake timeout")
	nameserverAddr      = flag.String("nameserver-addr", "8.8.8.8", "set up default fallback address")
	connectionTimeout   = flag.Duration("connection-timeout", 0, "set connection timeout")
	tcpKeepIdle         = flag.Duration("tcp-keep-idle", 1*time.Second, "set TCP keepalive idle in seconds")
	tcpKeepInterval     = flag.Duration("tcp-keep-interval", 1*time.Second, "set TCP keepalive interval, in seconds")
	tcpKeepFailAfter    = flag.Duration("tcp-keep-failafter", 3*time.Second, "set TCP keepalive fail after value")
	readTimeout         = flag.Duration("read-timeout", 60*time.Second, "set read timeout if OS do not support TCP Keep-alive, in seconds")
	writeTimeout        = flag.Duration("write-timeout", 60*time.Second, "set write timeout if OS do not support TCP Keep-alive, in seconds")

	filename = flag.String("filename", "inputs.json", "file which content will be sent with POST request")
	endpoint = flag.String("endpoint", "https://httpbin.org/post", "endpoint where data will be sent")
)

func main() {
	flag.Parse()

	fmt.Println("DROBOT")

	clientConfig := client.Config{
		ResolveTimeout:      *resolveTimeout,
		TLSHandshakeTimeout: *tlsHandshakeTimeout,
		NameserverAddr:      *nameserverAddr,
		ConnectionTimeout:   *connectionTimeout,
		TCPKeepIdleTime:     *tcpKeepIdle,
		TCPKeepIntervalTime: *tcpKeepInterval,
		TCPKeepFailAfter:    *tcpKeepFailAfter,
		ReadTimeout:         *readTimeout,
		WriteTimeout:        *writeTimeout,
	}

	client, err := client.NewClient(clientConfig)
	if err != nil {
		panic(err)
	}

	contents, err := ioutil.ReadFile(*filename)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(http.MethodPost, *endpoint, bytes.NewBuffer(contents))
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
