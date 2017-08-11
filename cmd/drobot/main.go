package main

import (
	"flag"
)

// flags
var (
	resolveTimeout      = flag.Int("resolve-timeout", 0, "set resolve timeout")
	tlsHandshakeTimeout = flag.Int("tls-handshake-timeout", 0, "set TLS handshake timeout")
	fallbackAddr        = flag.String("fallback-addr", "8.8.8.8", "set up default fallback address")
	connectionTimeout   = flag.Int("connection-timeout", 0, "set connection timeout")
	tcpKAIdle           = flag.Int("tcp-keepalive-idle", 1, "set TCP keepalive idle in seconds")
	tcpKAInterval       = flag.Int("tcp-keepalive-interval", 1, "set TCP keepalive interval, in seconds")
	tcpKAFailAfter      = flag.Int("tcp-keepalive-failafter", 3, "set TCP keepalive fail after value")
	rwTimeout           = flag.Int("rw-timeout", 60, "set read/write timeout if OS do not support TCP Keep-alive, in seconds")
)

func main() {
	flag.Parse()

	fmt.Println("DROBOT")
}
