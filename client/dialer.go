package client

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

// Dialer is a custom dialer with a DialContext method
type Dialer struct {
	*net.Dialer

	keepIdleTime     int
	keepIntervalTime int
	keepFailAfter    int
	nameserverAddr   string

	// read/write timeout is set if kernel do not support TCP keep-alive
	readTimeout  int
	writeTimeout int
}

func (d *Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	var host, port string
	var err error

	addrs := []string{}

	host, port, err = net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	// if host is already address
	// then skip DNS address resolution
	if ip := net.ParseIP(host); ip != nil {
		addrs = append(addrs, host)
	} else {
		addrs, err = getHostAddrs(ctx, host, d.nameserverAddr)
		if err != nil {
			return nil, err
		}
	}

	var prevAddr string
	var conn net.Conn
	for i := 0; i < 2; i++ {
		// if it is second try and address slices are only one
		// break, second try will not make it successful
		if len(addrs) < 2 && i > 0 {
			break
		}

		addr := getRandomAddr(addrs)

		// if prev addr and current is the same
		// then try one more time
		if prevAddr == addr {
			addr = getRandomAddr(addrs)
		}

		conn, err = d.Dialer.DialContext(ctx, network, addr+":"+port)
		if err != nil {
			continue
		} else {
			break
		}
	}

	// all tries are failed
	// return error
	if err != nil {
		return nil, err
	}

	err = setTCPKeepAlive(conn, d.keepIdleTime, d.keepIntervalTime, d.keepFailAfter)
	if err != nil {
		err = conn.SetReadDeadline(time.Now().Add(time.Duration(d.readTimeout) * time.Second))
		if err != nil {
			return conn, err
		}

		err = conn.SetWriteDeadline(time.Now().Add(time.Duration(d.writeTimeout) * time.Second))
		if err != nil {
			return conn, err
		}
	}

	return conn, nil
}

func setTCPKeepAlive(c net.Conn, idleTime, intervalTime, count int) error {
	conn, ok := c.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("invalid connection type: %T", c)
	}

	file, err := conn.File()
	if err != nil {
		return err
	}

	fd := int(file.Fd())

	err = setSockKeepIdleTime(fd, idleTime)
	if err != nil {
		return err
	}

	err = setSockKeepIntervalTime(fd, intervalTime)
	if err != nil {
		return err
	}

	err = setSockKeepCount(fd, count)
	if err != nil {
		return err
	}

	return nil
}

// seconds returns seconds of Duration
func seconds(d time.Duration) int {
	d += (time.Second - time.Nanosecond)
	return int(d.Seconds())
}

func getHostAddrs(ctx context.Context, host, ns string) ([]string, error) {
	addrs, err := net.DefaultResolver.LookupHost(ctx, host)
	if err == nil {
		return addrs, nil
	}

	// if address is not found using local resolver
	// get address from remote name server
	fmt.Printf("failed to resolve host '%s' with local resolver, searching in %s", host, ns)

	client := new(dns.Client)
	msg := new(dns.Msg)

	msg.SetQuestion(dns.Fqdn(host), dns.TypeA)

	reply, _, err := client.Exchange(msg, ns)
	if err != nil {
		return nil, err
	}

	if reply.Rcode != dns.RcodeSuccess {
		return nil, nil
	}

	if len(reply.Answer) < 1 {
		return nil, fmt.Errorf("can't resolv domain: %s", host)
	}

	// get IP addresses
	for _, a := range reply.Answer {
		addrs = append(addrs, dns.Field(a, 1))
	}

	return addrs, nil
}
