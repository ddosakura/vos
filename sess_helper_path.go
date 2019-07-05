package vos

import (
	"io/ioutil"
	"path"
)

func (s *Session) pathHelper(line string) []string {
	// TODO: path_helper
	names := make([]string, 0)
	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names
}

// Ls support relative path
func (s *Session) Ls(p string, c *LsConfig) ([]string, error) {
	p = path.Join(s.pwd, p)
	return s.os.Ls(p, c)
}
