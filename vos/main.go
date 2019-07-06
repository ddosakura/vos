package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ddosakura/vos"
	"github.com/manifoldco/promptui"
	"github.com/spf13/afero"
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

func auth(v *vos.OS, d interface{}) (*vos.Session, error) {
	ds := d.([2]string)
	user := ds[0]
	pass := ds[1]
	if (user == "root" && pass == "pass") || (user == "sakura" && pass == "123456") {
		s := v.NewSession(user)
		s.In = os.Stdin
		s.Out = os.Stdout
		return s, nil
	}
	return nil, ErrAuthFail
}

func main() {
	v := vos.New()
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		v.LogMode = vos.LmD
	}
	v.Init = func(api *vos.OS) ([]string, error) {
		// 挂载文件系统
		fs := afero.NewOsFs()
		api.Mount(&vos.Fs{
			Type:  "/dev/sda2",
			Point: "/",
			Mount: afero.NewBasePathFs(fs, "./sda2"),
		})
		api.Mount(&vos.Fs{
			Type:  "/dev/sda1",
			Point: "/boot",
			Mount: afero.NewBasePathFs(fs, "./sda1"),
		})
		api.Mount(&vos.Fs{
			Type:  "/dev/sdb1",
			Point: "/home",
			Mount: afero.NewBasePathFs(fs, "./sdb1"),
		})

		return []string{}, nil
	}
	if e := v.Run(); e != nil {
		panic(e)
	}

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
		if s, err = auth(v, [2]string{u, p}); err == nil {
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
