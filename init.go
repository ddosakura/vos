package vos

import (
	"github.com/spf13/afero"
)

func (v *OS) defaultInit() (afero.Fs, []string, error) {
	return nil, nil, ErrInitScriptNotFound
}

// Hostname of OS
func (v *OS) Hostname() string {
	return "OS"
}
