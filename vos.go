package vos

import (
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/afero"
)

// raw
const (
	Version = "v0.0.1"
)

// OS Virtual Operating System
type OS struct {
	// OS Core
	Fs          afero.Fs
	SyscallList []string
	Init        func()

	// Login
	Auth func(interface{}) (*Session, error)

	// Controller
	api     chan Syscall
	stop    chan interface{}
	session map[string]*Session
	sess    chan sess
}

// Syscall of VOS
type Syscall struct {
	Name    string
	Payload interface{}
}

type sess struct {
	uuid string
	s    *Session
}

// New VOS
func New() *OS {
	v := &OS{
		Fs:          nil,
		SyscallList: []string{},
		Init:        nil,
		Auth:        nil,

		api:     make(chan Syscall),
		stop:    make(chan interface{}),
		session: make(map[string]*Session),
		sess:    make(chan sess),
	}
	return v
}

// Run VOS
func (v *OS) Run() {
	if v.Init == nil {
		v.Init = v.defaultInit
	}

	for {
		select {
		case <-v.stop:
			return
		case s := <-v.sess:
			if s.s == nil {
				delete(v.session, s.uuid)
				s.s.os = nil
			} else {
				UUID := uuid.NewV4().String()
				v.session[UUID] = s.s
				s.s.uuid = UUID
			}
		}

	}
}

// Syscall for Handler
func (v *OS) Syscall() <-chan Syscall {
	return v.api
}

// Stop VOS
func (v *OS) Stop() {
	v.stop <- nil
}
