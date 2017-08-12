package client

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Dialer is a custom dialer with a DialContext method
type Dialer struct {
	*net.Dialer

	keepIdleTime     int
	keepIntervalTime int
	keepFailAfter    int

	// read/write timeout is set if kernel do not support TCP keep-alive
	readTimeout  int
	writeTimeout int
}

func (d *Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	conn, err := d.Dialer.DialContext(ctx, network, address)
	if err != nil {
		return conn, err
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
