package main

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/chzyer/readline"
	"github.com/ddosakura/pbqp"
	"github.com/ddosakura/vos/proto/auth"
	"github.com/manifoldco/promptui"
)

// msg(s)
const (
	MsgAuthFail = "Permission denied, please try again."
)

// error(s)
var (
	ErrNoInput = errors.New("No Input, please input again")
)

var (
	promptUser = promptui.Prompt{
		Label:    "username",
		Validate: authValidate,
	}
	promptPass = promptui.Prompt{
		Label:    "password",
		Validate: authValidate,
		Mask:     '*',
	}
)

func authValidate(input string) error {
	if input == "" {
		return ErrNoInput
	}
	return nil
}

func main() {
	conn, err := net.Dial("unix", os.Args[1])
	if err != nil {
		panic(err)
	}

	for {
		var err error
		var u, p string
		if u, err = promptUser.Run(); err != nil {
			panic(err)
		}
		if p, err = promptPass.Run(); err != nil {
			panic(err)
		}
		if err := pbqp.Write(conn, &auth.Auth{
			Ver:  1,
			Type: auth.Type_Password,
			User: u,
			Pass: p,
		}); err != nil {
			panic(err)
		}
		r := new(auth.Result)
		if err := pbqp.Read(conn, r); err != nil {
			panic(err)
		}
		if r.Pass {
			fmt.Println(r.Welcome)
			break
		}
		fmt.Printf("\n%s\n\n", MsgAuthFail)
	}

	cli, err := readline.NewRemoteCli(conn)
	if err != nil {
		panic(err)
	}
	if err = cli.Serve(); err != nil {
		panic(err)
	}
}
