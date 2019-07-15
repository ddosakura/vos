package vos

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

// https://github.com/chzyer/readline/blob/master/example/readline-demo/readline-demo.go

func rlCFGBuilder(s *Session) *readline.Config {
	return &readline.Config{
		Prompt:          Prompt,
		HistoryFile:     s.historyFile,
		AutoComplete:    s.autoCompleter(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	}
}

func rlHandlerBuilder(s *Session) func(l *readline.Instance) {
	return func(l *readline.Instance) {
		s.Info("Session Key-Ring Initializing ...")
		//setPasswordCfg := KR(l)

		var lastError error
		for {
			if PromptTpl != nil {
				fmt.Fprintln(l.Stdout(), s.prompt(lastError))
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
				return
			//case strings.HasPrefix(line, "echo"):
			//	line := strings.TrimSpace(line[4:])
			//	_, lastError = s.println(line)
			//case strings.HasPrefix(line, "uname"):
			//	if strings.TrimSpace(line[5:]) == "-a" {
			//		_, lastError = s.println("VOS %s %s %s %s DDoSakura/VOS", s.os.Hostname(), s.os.Version(), s.os.BuildTime(), BitIf())
			//	} else {
			//		_, lastError = s.println("VOS")
			//	}
			//case strings.HasPrefix(line, "pwd"):
			//	if line == "pwd" {
			//		_, lastError = s.println(s.pwd)
			//	} else {
			//		_, lastError = s.println("pwd: too many arguments")
			//		if lastError == nil {
			//			lastError = errors.New("C:1")
			//		}
			//	}
			//case strings.HasPrefix(line, "cd"):
			//	p0 := strings.TrimSpace(line[2:])
			//	p := p0
			//	if p == "" {
			//		s.pwd = home(s.user)
			//		continue
			//	}
			//	if p[0] == '~' {
			//		p = strings.Replace(p, "~", home(s.user), 1)
			//	} else if p[0] != '/' {
			//		p = path.Join(s.pwd, p)
			//	}
			//	fs, _ := s.os.FindFile(p)
			//	if fs == nil {
			//		_, lastError = s.println("cd: 没有那个文件或目录: %s", p0)
			//		if lastError == nil {
			//			lastError = errors.New("C:1")
			//		}
			//		continue
			//	}
			//	s.pwd = p
			//	_, lastError = s.println("")
			//case strings.HasPrefix(line, "ls"):
			//	c := new(LsConfig)
			//	c.Color = true
			//	// TODO: more config & ls other path
			//	if strings.TrimSpace(line[2:]) == "-a" {
			//		c.A = true
			//	}
			//	ls, err := s.os.Ls(s.pwd, c)
			//	if err != nil {
			//		lastError = err
			//		continue
			//	}
			//	// TODO: fmt
			//	res := ""
			//	for _, o := range ls {
			//		res += o + "	"
			//	}
			//	_, lastError = s.println(res)

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
			//case line == "setpassword":
			//	pswd, err := l.ReadPasswordWithConfig(setPasswordCfg)
			//	if err == nil {
			//		println("you set:", strconv.Quote(string(pswd)))
			//	}
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
				return
			case line == "sleep":
				s.Info("sleep 4 second")
				time.Sleep(4 * time.Second)
			case line == "":
			default:
				s.Info("you said:", strconv.Quote(line))
			}
		}
	}
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
		readline.PcItem("uname", readline.PcItem("-a")),
		// TODO: readline.PcItem("shutdown"),
		readline.PcItem("pwd"),
		//readline.PcItem("cd", readline.PcItemDynamic(s.pathHelper)),
		//readline.PcItem("ls",
		//	readline.PcItem("-a",
		//		readline.PcItemDynamic(s.pathHelper),
		//	),
		//	readline.PcItemDynamic(s.pathHelper),
		//),
		//// TODO: commands
		//readline.PcItem("df", readline.PcItemDynamic(s.pathHelper)),

		//readline.PcItem("touch", readline.PcItemDynamic(s.pathHelper)),
		//readline.PcItem("mkdir", readline.PcItemDynamic(s.pathHelper)),
		//readline.PcItem("rm", readline.PcItemDynamic(s.pathHelper)),
		//readline.PcItem("cat", readline.PcItemDynamic(s.pathHelper)),
		//readline.PcItem("tail", readline.PcItemDynamic(s.pathHelper)),

		//readline.PcItem("mount", readline.PcItemDynamic(s.pathHelper)),
		//readline.PcItem("umount", readline.PcItemDynamic(s.pathHelper)),

		// TODO: 环境变量
		//readline.PcItem("export", readline.PcItemDynamic(s.pathHelper)),

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
