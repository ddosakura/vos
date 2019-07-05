package vos

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"
)

// Prompt
var (
	defaultPromptTpl, _ = template.New("Prompt").Parse(fmt.Sprintf("\n%s {{ .User }} @ {{ .Host }} in {{ .PWD }} [{{ .Time }}] {{ .Error }}", Blue("#")))

	defaultPrompt = Red("$") + " "
)

func (s *Session) prompt(err error) string {
	buf := new(bytes.Buffer)
	errText := ""
	if err != nil {
		errText = err.Error()
	}
	s.PromptTpl.Execute(buf, struct {
		User  string
		Host  string
		PWD   string
		Time  string
		Error string
	}{
		User:  Blue(s.user),
		Host:  Green(s.os.Hostname()),
		PWD:   Yellow(s.showPWD()),
		Time:  time.Now().Format("15:04:05"), // 2006-01-02 15:04:05
		Error: Red(errText),
	})
	return buf.String()
}

func (s *Session) showPWD() string {
	h := home(s.user)
	if strings.HasPrefix(s.pwd, h) {
		return strings.Replace(s.pwd, h, "~", 1)
	}
	return s.pwd
}

func (s *Session) print(f string, v ...interface{}) (n int, e error) {
	return s.Out.Write([]byte(fmt.Sprintf(f, v...)))
}

func (s *Session) println(f string, v ...interface{}) (n int, e error) {
	return s.Out.Write([]byte(fmt.Sprintf(f+"\n", v...)))
}
