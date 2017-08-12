package client

import (
	"os"
	"syscall"
)

// from netinet/tcp.h (OS X 10.9.4)
const (
	_TCP_KEEPINTVL = 0x101 /* interval between keepalives */
	_TCP_KEEPCNT   = 0x102 /* number of keepalives before close */
)

func setSockKeepIdleTime(fd, idleTime int) error {
	err := os.NewSyscallError("setsockopt", syscall.SetsockoptInt(
		fd,
		syscall.IPPROTO_TCP,
		syscall.TCP_KEEPALIVE,
		idleTime,
	))

	return err
}

func setSockKeepIntervalTime(fd, intervalTime int) error {
	err := os.NewSyscallError("setsockopt", syscall.SetsockoptInt(
		fd,
		syscall.IPPROTO_TCP,
		_TCP_KEEPINTVL,
		intervalTime,
	))

	return err
}

func setSockKeepCount(fd, intervalTime int) error {
	err := os.NewSyscallError("setsockopt", syscall.SetsockoptInt(
		fd,
		syscall.IPPROTO_TCP,
		_TCP_KEEPCNT,
		intervalTime,
	))

	return err
}
