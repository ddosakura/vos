package vos

import (
	"fmt"
	"log"
	"os"
	"sync"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/afero"
)

// OS Virtual Operating System
type OS struct {
	Issuer  string
	LogMode LogMode
	Logger  *log.Logger

	// Init & OS Core
	Init        func(*OS) ([]string, error)
	fs          map[string]*Fs // point => fs
	fsLock      *sync.RWMutex
	syscallList []string

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
		Issuer:  "core",
		LogMode: LmI,
		Logger:  log.New(os.Stderr, "[VOS] ", log.LstdFlags),

		Init:   nil,
		fs:     make(map[string]*Fs),
		fsLock: new(sync.RWMutex),
		//syscallList: []string{},

		api:     make(chan Syscall),
		stop:    make(chan interface{}),
		session: make(map[string]*Session),
		sess:    make(chan sess),
	}
	return v
}

// Run VOS
func (v *OS) Run() error {
	v.Info("OS Initializing ...")
	// Init
	if v.Init == nil {
		v.Init = v.defaultInit
	}
	s, e := v.Init(v)
	if e != nil {
		return e
	}
	v.syscallList = s
	if e = v.Mount(&Fs{
		Type:  "buildin",
		Point: "/sys",
		// TODO: 替换成内建 fs
		Mount: afero.NewBasePathFs(afero.NewOsFs(), "./buildin"),
	}); e != nil {
		return e
	}
	v.syscallList = s
	if e = v.Mount(&Fs{
		Type:  "tmpfs",
		Point: "/tmp",
		Mount: afero.NewMemMapFs(),
	}); e != nil {
		return e
	}
	if e = v.Mount(&Fs{
		Type:  "",
		Point: "/run",
		Mount: afero.NewMemMapFs(),
	}); e != nil {
		return e
	}

	go func() {
		for {
			select {
			case <-v.stop:
				return
			case s := <-v.sess:
				// TODO: 应当触发一个系统事件
				if s.s == nil {
					v.Info(fmt.Sprintf("Session[UUID=%s] Close", s.uuid))
					delete(v.session, s.uuid)
					s.s.os = nil
				} else {
					UUID := uuid.NewV4().String()
					s.uuid = UUID
					v.Info(fmt.Sprintf("Session[UUID=%s] Created", s.uuid))
					v.session[UUID] = s.s
					s.s.uuid = UUID
				}
			}

		}
	}()

	v.Info("Executing initialization Script ...")
	// TODO: 执行初始化脚本 /boot/init

	return nil
}

// Syscall for Handler
func (v *OS) Syscall() <-chan Syscall {
	return v.api
}

// Stop VOS
func (v *OS) Stop() {
	v.Info("OS Shutdown ...")
	v.stop <- nil
}
