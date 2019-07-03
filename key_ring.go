package vos

import (
	"fmt"

	"github.com/chzyer/readline"
)

// TODO: Key Ring

// KR Temp
func KR(l *readline.Instance) *readline.Config {
	setPasswordCfg := l.GenPasswordConfig()
	setPasswordCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		l.SetPrompt(fmt.Sprintf("Enter password(%v): ", len(line)))
		l.Refresh()
		return nil, 0, false
	})
	return setPasswordCfg
}
