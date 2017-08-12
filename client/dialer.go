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
	var IPAddr, IPAddrSecondary string
	var host, port string
	var err error

	// if host is already address
	// then skip DNS address resolution
	if ip := net.ParseIP(address); ip != nil {
		IPAddr = address
	} else {
		host, port, err = net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}

		addrs, err := net.DefaultResolver.LookupHost(ctx, host)
		if err != nil {
			return nil, err
		}

		// if address is not found using local resolver
		// get address from remote name server
		if err != nil || len(addrs) == 0 {
			client := new(dns.Client)
			msg := new(dns.Msg)

			msg.SetQuestion(dns.Fqdn(address), dns.TypeA)

			reply, _, err := client.Exchange(msg, d.nameserverAddr)
			if err != nil {
				return nil, err
			}

			if reply.Rcode != dns.RcodeSuccess {
				return nil, nil
			}

			for _, a := range reply.Answer {
				// first field is IP
				IPAddr = dns.Field(a, 1)
				break
			}
		} else {
			// if we have multiple addresses
			// choose one randomly
			if len(addrs) > 1 {
				IPAddr = getRandomAddr(addrs)
				IPAddrSecondary = getRandomAddr(addrs)
			} else {
				IPAddr = addrs[0]
			}

		}
	}

	println(IPAddr, IPAddrSecondary, port)
	conn, err := d.Dialer.DialContext(ctx, network, IPAddr+":"+port)
	if err != nil {
		conn, err = d.Dialer.DialContext(ctx, network, IPAddrSecondary+":"+port)
		if err != nil {
			return conn, err
		}
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

		return conn, err
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
