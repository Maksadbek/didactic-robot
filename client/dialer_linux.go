package client

import (
	"os"
	"syscall"
)

func setSockKeepIdleTime(fd, secs int) error {
	err := os.NewSyscallError("setsockopt", syscall.SetsockoptInt(
		fd,
		syscall.IPPROTO_TCP,
		syscall.TCP_KEEPIDLE,
		secs,
	))

	return err
}

func setSockKeepIntervalTime(fd, secs int) error {
	err := os.NewSyscallError("setsockopt", syscall.SetsockoptInt(
		fd,
		syscall.IPPROTO_TCP,
		syscall.TCP_KEEPINTVL,
		secs,
	))

	return err
}

func setSockKeepCount(fd, secs int) error {
	err := os.NewSyscallError("setsockopt", syscall.SetsockoptInt(
		fd,
		syscall.IPPROTO_TCP,
		syscall.TCP_KEEPCNT,
		secs,
	))

	return err
}
