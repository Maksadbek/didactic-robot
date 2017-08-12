package client

func setSockKeepIdleTime(fd, secs int) error {
	err := os.NewSyscallError("setsockopt", syscall.SetsockoptInt(
		fd,
		syscall.IPPROTO_TCP,
		syscall.TCP_KEEPALIVE,
		seconds(idleTime),
	))

	return err
}

func setSockKeepIntervalTime(fd, secs int) error {
	err := os.NewSyscallError("setsockopt", syscall.SetsockoptInt(
		fd,
		syscall.IPPROTO_TCP,
		syscall.TCP_KEEPINTVL,
		seconds(intervalTime),
	))

	return err
}

func setSockKeepCount(fd, secs int) error {
	err := os.NewSyscallError("setsockopt", syscall.SetsockoptInt(
		fd,
		syscall.IPPROTO_TCP,
		syscall.TCP_KEEPCount,
		seconds(intervalTime),
	))

	return err
}
