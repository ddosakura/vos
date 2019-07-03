package vos

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/chzyer/readline"
)

// Session of VOS
type Session struct {
	os   *OS
	pwd  string
	uuid string

	PromptTpl   *template.Template
	Prompt      string
	HistoryFile string
	LogFile     string

	In   io.ReadCloser
	Out  io.WriteCloser
	User string
}

// NewSession for VOS
func (v *OS) NewSession() *Session {
	s := &Session{
		os:  v,
		pwd: "~",
	}
	v.sess <- sess{s: s}
	s.PromptTpl = defaultPromptTpl
	s.Prompt = defaultPrompt
	for s.uuid == "" {
		time.Sleep(time.Millisecond * 100)
	}
	s.HistoryFile = "./" + s.uuid + "_history"
	s.LogFile = "./" + s.uuid + ".log"
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
	// https://github.com/chzyer/readline/blob/master/example/readline-demo/readline-demo.go
	l, err := readline.NewEx(&readline.Config{
		Prompt:          s.Prompt,
		HistoryFile:     s.HistoryFile,
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		return err
	}
	defer l.Close()

	setPasswordCfg := KR(l)

	var lastError error
	// TODO: log.SetOutput(...)
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
			usage(l.Stderr())
		case line == "setpassword":
			pswd, err := l.ReadPasswordWithConfig(setPasswordCfg)
			if err == nil {
				println("you set:", strconv.Quote(string(pswd)))
			}
		case strings.HasPrefix(line, "setprompt"):
			if len(line) <= 10 {
				log.Println("setprompt <prompt>")
				break
			}
			l.SetPrompt(line[10:])
		case strings.HasPrefix(line, "say"):
			line := strings.TrimSpace(line[3:])
			if len(line) == 0 {
				log.Println("say what?")
				break
			}
			go func() {
				for range time.Tick(time.Second) {
					log.Println(line)
				}
			}()
		case line == "bye":
			goto exit
		case line == "sleep":
			log.Println("sleep 4 second")
			time.Sleep(4 * time.Second)
		case line == "":
		default:
			log.Println("you said:", strconv.Quote(line))
		}
	}
exit:
	return nil
}

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
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

var completer = readline.NewPrefixCompleter(
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

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
