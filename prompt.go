package vos

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	glog "github.com/ddosakura/ghost-log"
)

// Prompt
var (
	PromptTpl, _ = template.New("Prompt").Parse(fmt.Sprintf("\n%s {{ .User }} @ {{ .Host }} in {{ .PWD }} [{{ .Time }}] {{ .Error }}", glog.Blue("#")))
	Prompt       = glog.Red("$") + " "
)

func (s *Session) prompt(err error) string {
	buf := new(bytes.Buffer)
	errText := ""
	if err != nil {
		errText = err.Error()
	}
	PromptTpl.Execute(buf, struct {
		User  string
		Host  string
		PWD   string
		Time  string
		Error string
	}{
		User:  glog.Blue(s.user),
		Host:  glog.Green(s.os.Syscall(SignHostname)),
		PWD:   glog.Yellow(s.pwd),
		Time:  time.Now().Format("15:04:05"), // 2006-01-02 15:04:05
		Error: glog.Red(errText),
	})
	return buf.String()
}
