package vos

import (
	"os"

	"github.com/spf13/afero"
)

// Fs of VOS
type Fs struct {
	Type  string
	Point string
	Mount afero.Fs
}

// Mount Fs
func (v *OS) Mount(fs *Fs) error {
	v.Info("Mount", fs.Point)
	v.fsLock.Lock()
	defer v.fsLock.Unlock()
	p := fs.Point
	if p == "" {
		return ErrMountPointNotFound
	}
	if p == "/" {
		if v.fs["/"] != nil {
			return os.ErrExist
		}
	} else {
		if v.fs["/"] == nil {
			return ErrMountPointNotFound
		}
		if fs, _ := v.findFileObject(p); fs != nil {
			return os.ErrExist
		}
	}

	v.fs[p] = fs
	return nil
}

// UMount Fs
func (v *OS) UMount(point string) error {
	v.Info("UMount", point)
	v.fsLock.Lock()
	defer v.fsLock.Unlock()
	delete(v.fs, point)
	return nil
}

// FindFile wrap findFileObject with lock
func (v *OS) FindFile(path string) (afero.Fs, string) {
	v.fsLock.RLock()
	defer v.fsLock.RUnlock()
	return v.findFileObject(path)
}

// Hostname of OS
func (v *OS) Hostname() string {
	v.Debug("Call v.Hostname()")
	fs, path := v.FindFile("/proc/sys/kernel/hostname")
	bs, err := afero.ReadFile(fs, path)
	if err != nil {
		return "VOS"
	}
	return string(bs)
}

// LsConfig for `ls`
type LsConfig struct {
	A     bool // include .*
	Color bool
}

// Ls - `ls`
func (v *OS) Ls(path string, c *LsConfig) ([]string, error) {
	fs, path := v.FindFile(path)
	if fs == nil {
		return nil, os.ErrNotExist
	}
	fis, err := afero.ReadDir(fs, path)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(fis)+2)
	if c.A {
		if c.Color {
			res = append(res, Blue("."), Blue(".."))
		} else {
			res = append(res, ".", "..")
		}
	}
	for _, fi := range fis {
		if fi.Name()[0] == '.' && !c.A {
			continue
		}

		if c.Color {
			if fi.IsDir() {
				// Directory
				res = append(res, Blue(fi.Name()))
			} else if fi.Mode().Perm()&PermX == PermX {
				// Executable File
				res = append(res, Green(fi.Name()))
			} else {
				// Normal File
				res = append(res, fi.Name())
			}
		} else {
			res = append(res, fi.Name())
		}
	}
	return res, nil
}
