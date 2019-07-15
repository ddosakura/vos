package vos

import (
	"fmt"
	"net"
	"sync"

	"github.com/chzyer/readline"
	glog "github.com/ddosakura/ghost-log"
	"github.com/ddosakura/pbqp"
	"github.com/ddosakura/vos/proto/auth"
	uuid "github.com/satori/go.uuid"
)

// Ext of VOS
type Ext interface {
	Init(*Base) error
	Syscall(args ...string) interface{}
}

// Base VOS
type Base struct {
	*glog.Log

	exts []Ext

	listener net.Listener
	sessMap  map[string]*Session
	sessLock *sync.RWMutex

	// TODO: 虚拟进程
}

// New VOS
func New(exts ...Ext) *Base {
	v := &Base{
		Log:      glog.New(),
		sessMap:  make(map[string]*Session),
		sessLock: new(sync.RWMutex),
		exts:     exts,
	}
	return v
}

// Run Base-VOS
func (b *Base) Run(l net.Listener) {
	b.Info("OS Initializing ...")
	for _, m := range b.exts {
		if err := m.Init(b); err != nil {
			b.Error(err)
		}
	}
	b.Info("OS initialization completed!")
	b.listener = l
	for {
		conn, err := b.listener.Accept()
		go func() {
			defer glog.StopPanic()
			if err != nil {
				b.Error(err)
			}
			num, ok := b.Syscall(SignLoginRetryTimes).(int)
			if !ok {
				b.Error(ErrSyscallError)
			}
			var user string
			for num > 0 {
				r := new(auth.Auth)
				if err := pbqp.Read(conn, r); err != nil {
					b.Warn(err)
				}
				res := new(auth.Result)
				res.Pass = false
				switch r.Type {
				case auth.Type_Password:
					// TODO: more users
					if r.User == "root" && r.Pass == "123456" {
						res.Pass = true
						// TODO: welcome
						res.Welcome = `Success`
						user = r.User
					}
				case auth.Type_PublicKey:
					// TODO:
					b.Warn(ErrAuthProtocol)
					// TODO: if success, then user = "xxx"
				}
				if err := pbqp.Write(conn, res); err != nil {
					b.Warn(err)
				}
				if res.Pass {
					break
				}
				num--
			}
			if num == 0 {
				conn.Close()
				return
			}

			s := b.NewSession(func(s *Session) {
				s.user = user
				s.uuid = uuid.NewV4().String()
				b.Info(fmt.Sprintf("Session[UUID=%s] Created", s.uuid))
			})
			cfg := rlCFGBuilder(s)
			fn := rlHandlerBuilder(s)
			rl, err := readline.HandleConn(*cfg, conn)
			if err != nil {
				conn.Close()
				b.Error(err)
			}
			fn(rl)
			conn.Close()
		}()
	}
}
