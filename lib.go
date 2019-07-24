package vos

import (
	"fmt"
	"log"
	"net"
	"os"
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

	sessMap  map[string]*Session
	sessLock *sync.RWMutex

	Auth func(*auth.Auth, *auth.Result) (string, error)

	// TODO: 虚拟进程
}

// New VOS
func New(exts ...Ext) *Base {
	v := &Base{
		Log:      glog.New(),
		exts:     exts,
		sessMap:  make(map[string]*Session),
		sessLock: new(sync.RWMutex),
		Auth:     defaultAuth,
	}
	v.Log.Logger = log.New(os.Stderr, "[VOS] ", log.LstdFlags)
	return v
}

func defaultAuth(req *auth.Auth, res *auth.Result) (string, error) {
	switch req.Type {
	case auth.Type_Password:
		if req.User == "root" && req.Pass == "123456" {
			res.Pass = true
			res.Welcome = `Success`
			return req.User, nil
		}
	case auth.Type_PublicKey:
		return "", ErrAuthProtocol
	default:
		return "", ErrAuthProtocol
	}
	return "", nil
}

// Start Base-VOS
func (b *Base) Start() error {
	b.Info("OS Initializing ...")
	for _, m := range b.exts {
		if err := m.Init(b); err != nil {
			return err
		}
	}
	b.Info("OS initialization completed!")
	return nil
}

func (b *Base) Accept(conn net.Conn) {
	defer glog.StopPanic()
	num, ok := b.Syscall(SignLoginRetryTimes).(int)
	if !ok {
		b.Error(ErrSyscallError)
	}
	var user string
	for num > 0 {
		r := new(auth.Auth)
		if err := pbqp.Read(conn, r); err != nil {
			b.Error(err)
		}
		res := new(auth.Result)
		res.Pass = false
		if u, err := b.Auth(r, res); err != nil {
			b.Warn(err)
		} else if res.Pass {
			user = u
		}
		if err := pbqp.Write(conn, res); err != nil {
			b.Error(err)
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
}
