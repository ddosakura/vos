package vos

import (
	"log"
	"os"

	glog "github.com/ddosakura/ghost-log"
)

// Session of VOS
type Session struct {
	*glog.Log
	os          *Base
	user        string
	uuid        string
	historyFile string

	pwd string
}

// NewSession for Base-VOS
func (b *Base) NewSession(fn func(*Session)) *Session {
	s := &Session{
		os: b,

		pwd: "~",
	}
	fn(s)
	s.Log = &glog.Log{
		Mode:   b.Mode,
		Logger: log.New(os.Stderr, "[VOS-Session "+s.uuid+"("+s.user+")] ", log.LstdFlags),
	}
	s.historyFile = "./" + s.user + "_history"
	return s
}
