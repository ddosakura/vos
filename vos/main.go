package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ddosakura/vos"
	"github.com/manifoldco/promptui"
)

// error(s)
var (
	ErrAuthFail = errors.New("Auth Fail")
)

var (
	promptUser = promptui.Prompt{
		Label:    "username ",
		Validate: authValidate,
	}
	promptPass = promptui.Prompt{
		Label:    "password ",
		Validate: authValidate,
		Mask:     '*',
	}
)

func authValidate(input string) error {
	if input == "" {
		return errors.New("Please Input")
	}
	return nil
}

func main() {
	v := vos.New()
	v.Auth = func(d interface{}) (*vos.Session, error) {
		ds := d.([2]string)
		user := ds[0]
		pass := ds[1]
		if user == "root" && pass == "pass" {
			s := v.NewSession()
			s.In = os.Stdin
			s.Out = os.Stdout
			s.User = user
			return s, nil
		}
		return nil, ErrAuthFail
	}

	go v.Run()

	var s *vos.Session
	num := 3
	for num > 0 {
		var err error
		var u, p string
		if u, err = promptUser.Run(); err != nil {
			panic(err)
		}
		if p, err = promptPass.Run(); err != nil {
			panic(err)
		}
		if s, err = v.Auth([2]string{u, p}); err == nil {
			break
		} else if err != ErrAuthFail {
			panic(err)
		}
		fmt.Print("\nAuth Fail!\n\n")
		num--
	}
	if s != nil {
		s.Run()
	}
}
