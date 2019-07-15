package vos

import "errors"

// error(s)
var (
	ErrSyscallError = errors.New("Syscall Error")
	ErrAuthProtocol = errors.New("Auth Protocol Error")
)
