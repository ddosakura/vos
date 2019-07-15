package fs

import "github.com/spf13/afero"

// Fs of VOS
type Fs struct {
	Type  string
	Point string
	Mount afero.Fs
}
