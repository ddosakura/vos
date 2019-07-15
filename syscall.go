package vos

import (
	"strings"
)

// Syscall(s)
const (
	SignLoginRetryTimes = "s_login_retry_times"
	SignHostname        = "s_hostname"
)

// Syscall of Base-VOS
func (b *Base) Syscall(args ...string) interface{} {
	for _, m := range b.exts {
		r := m.Syscall(args...)
		if r != nil {
			return r
		}
	}

	switch strings.ToLower(args[0]) {
	case SignLoginRetryTimes:
		return 3
	case SignHostname:
		return "vos"
	}
	return nil
}
