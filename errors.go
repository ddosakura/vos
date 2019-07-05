package vos

import (
	"errors"
)

// error(s)
var (
	ErrInitScriptNotFound = errors.New("Init Script Not Found")
	ErrMountPointNotFound = errors.New("Mount Point Not Found")
)
