package vos

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/chzyer/readline"
)

// Session of VOS
type Session struct {
	DebugMode bool
	Logger    *log.Logger

	os   *OS
	user string
	pwd  string
	uuid string

	PromptTpl   *template.Template
	Prompt      string
	HistoryFile string

	In  io.ReadCloser
	Out io.WriteCloser
}

// NewSession for VOS
func (v *OS) NewSession(user string) *Session {
	s := &Session{
		DebugMode: v.DebugMode,
		Logger:    log.New(os.Stderr, "[VOS-Session] ", log.LstdFlags),
		os:        v,
		user:      user,
		pwd:       home(user),
	}
	v.sess <- sess{s: s}
	s.PromptTpl = defaultPromptTpl
	s.Prompt = defaultPrompt
	for s.uuid == "" {
		time.Sleep(time.Millisecond * 100)
	}
	s.HistoryFile = "./" + user + "_history"
	return s
}

// OS of Session
func (s *Session) OS() *OS {
	return s.os
}

// UUID of Session
func (s *Session) UUID() string {
	return s.uuid
}

// Exit Action
func (s *Session) Exit() {
	s.os.sess <- sess{uuid: s.uuid}
}

// Run Session
func (s *Session) Run() error {
	s.Info("Session Initializing ...")
	// https://github.com/chzyer/readline/blob/master/example/readline-demo/readline-demo.go
	l, err := readline.NewEx(&readline.Config{
		Prompt:          s.Prompt,
		HistoryFile:     s.HistoryFile,
		AutoComplete:    s.autoCompleter(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		return err
	}
	defer l.Close()

	s.Info("Session Key-Ring Initializing ...")
	setPasswordCfg := KR(l)

	var lastError error
	for {
		if s.PromptTpl != nil {
			fmt.Fprintln(s.Out, s.prompt(lastError))
		}
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "exit"):
			// TODO: exit code ?
			goto exit
		case strings.HasPrefix(line, "echo"):
			line := strings.TrimSpace(line[4:])
			_, lastError = s.println(line)
		case strings.HasPrefix(line, "pwd"):
			if line == "pwd" {
				_, lastError = s.println(s.pwd)
			} else {
				_, lastError = s.println("pwd: too many arguments")
				if lastError == nil {
					lastError = errors.New("C:1")
				}
			}
		case strings.HasPrefix(line, "cd"):
			p0 := strings.TrimSpace(line[2:])
			p := p0
			if p == "" {
				s.pwd = home(s.user)
				continue
			}
			if p[0] == '~' {
				p = strings.Replace(p, "~", home(s.user), 1)
			} else if p[0] != '/' {
				p = path.Join(s.pwd, p)
			}
			fs, _ := s.os.FindFile(p)
			if fs == nil {
				_, lastError = s.println("cd: 没有那个文件或目录: %s", p0)
				if lastError == nil {
					lastError = errors.New("C:1")
				}
				continue
			}
			s.pwd = p
			_, lastError = s.println("")
		case strings.HasPrefix(line, "ls"):
			c := new(LsConfig)
			c.Color = true
			// TODO: more config & ls other path
			if strings.TrimSpace(line[2:]) == "-a" {
				c.A = true
			}
			ls, err := s.os.Ls(s.pwd, c)
			if err != nil {
				lastError = err
				continue
			}
			// TODO: fmt
			res := ""
			for _, o := range ls {
				res += o + "	"
			}
			_, lastError = s.println(res)

		case strings.HasPrefix(line, "mode "):
			switch line[5:] {
			case "vi":
				l.SetVimMode(true)
			case "emacs":
				l.SetVimMode(false)
			default:
				println("invalid mode:", line[5:])
			}
		case line == "mode":
			if l.IsVimMode() {
				println("current mode: vim")
			} else {
				println("current mode: emacs")
			}
		case line == "login":
			pswd, err := l.ReadPassword("please enter your password: ")
			if err != nil {
				break
			}
			println("you enter:", strconv.Quote(string(pswd)))
		case line == "help":
			s.usage(l.Stderr())
		case line == "setpassword":
			pswd, err := l.ReadPasswordWithConfig(setPasswordCfg)
			if err == nil {
				println("you set:", strconv.Quote(string(pswd)))
			}
		case strings.HasPrefix(line, "setprompt"):
			if len(line) <= 10 {
				s.Info("setprompt <prompt>")
				break
			}
			l.SetPrompt(line[10:])
		case strings.HasPrefix(line, "say"):
			line := strings.TrimSpace(line[3:])
			if len(line) == 0 {
				s.Info("say what?")
				break
			}
			go func() {
				for range time.Tick(time.Second) {
					s.Info(line)
				}
			}()
		case line == "bye":
			goto exit
		case line == "sleep":
			s.Info("sleep 4 second")
			time.Sleep(4 * time.Second)
		case line == "":
		default:
			s.Info("you said:", strconv.Quote(line))
		}
	}
exit:
	return nil
}

func (s *Session) usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, s.autoCompleter().(*readline.PrefixCompleter).Tree("    "))
}

func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

func (s *Session) autoCompleter() readline.AutoCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("echo"),
		readline.PcItem("exit"),
		//readline.PcItem("shutdown"),
		readline.PcItem("pwd"),
		readline.PcItem("cd", readline.PcItemDynamic(s.pathHelper)),
		readline.PcItem("ls", readline.PcItemDynamic(s.pathHelper)),
		// TODO: commands
		readline.PcItem("df", readline.PcItemDynamic(s.pathHelper)),

		readline.PcItem("touch", readline.PcItemDynamic(s.pathHelper)),
		readline.PcItem("mkdir", readline.PcItemDynamic(s.pathHelper)),
		readline.PcItem("rm", readline.PcItemDynamic(s.pathHelper)),
		readline.PcItem("cat", readline.PcItemDynamic(s.pathHelper)),
		readline.PcItem("tail", readline.PcItemDynamic(s.pathHelper)),

		readline.PcItem("mount", readline.PcItemDynamic(s.pathHelper)),
		readline.PcItem("umount", readline.PcItemDynamic(s.pathHelper)),

		// TODO: 环境变量
		readline.PcItem("export", readline.PcItemDynamic(s.pathHelper)),

		readline.PcItem("mode",
			readline.PcItem("vi"),
			readline.PcItem("emacs"),
		),
		readline.PcItem("login"),
		readline.PcItem("say",
			readline.PcItemDynamic(listFiles("./"),
				readline.PcItem("with",
					readline.PcItem("following"),
					readline.PcItem("items"),
				),
			),
			readline.PcItem("hello"),
			readline.PcItem("bye"),
		),
		readline.PcItem("setprompt"),
		readline.PcItem("setpassword"),
		readline.PcItem("bye"),
		readline.PcItem("help"),
		readline.PcItem("go",
			readline.PcItem("build", readline.PcItem("-o"), readline.PcItem("-v")),
			readline.PcItem("install",
				readline.PcItem("-v"),
				readline.PcItem("-vv"),
				readline.PcItem("-vvv"),
			),
			readline.PcItem("test"),
		),
		readline.PcItem("sleep"),
	)
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
