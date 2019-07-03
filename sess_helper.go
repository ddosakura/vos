package vos

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/gookit/color"
)

// Prompt
var (
	Blue   = color.FgBlue.Render
	Green  = color.FgGreen.Render
	Yellow = color.FgYellow.Render
	Red    = color.FgRed.Render

	defaultPromptTpl, _ = template.New("Prompt").Parse(fmt.Sprintf("%s {{ .User }} @ {{ .Host }} in {{ .PWD }} [{{ .Time }}] {{ .Error }}", Blue("#")))

	defaultPrompt = Red("$") + " "
)

func (s *Session) prompt(err error) string {
	buf := new(bytes.Buffer)
	errText := ""
	if err != nil {
		errText = "E:" + err.Error()
	}
	s.PromptTpl.Execute(buf, struct {
		User  string
		Host  string
		PWD   string
		Time  string
		Error string
	}{
		User:  Blue(s.User),
		Host:  Green(s.os.Hostname()),
		PWD:   Yellow(s.pwd),
		Time:  time.Now().Format("15:04:05"), // 2006-01-02 15:04:05
		Error: Red(errText),
	})
	return buf.String()
}
