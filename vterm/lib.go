package vterm

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net"

	"github.com/ddosakura/pbqp"
	"github.com/ddosakura/vos"
	"github.com/ddosakura/vos/proto/auth"
	"github.com/manifoldco/promptui"
)

// msg(s)
const (
	MsgAuthFail = "Permission denied, please try again."
)

// error(s)
var (
	ErrNoInput  = errors.New("No Input, please input again")
	ErrAuthFail = errors.New(MsgAuthFail)
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

// Connect VOS
func Connect(conn net.Conn, pub *rsa.PublicKey, payload []byte, sig []byte) (err error) {
	if pub == nil {
		for {
			var u, p string
			if u, err = promptUser.Run(); err != nil {
				return
			}
			if p, err = promptPass.Run(); err != nil {
				return
			}
			if err = pbqp.Write(conn, &auth.Auth{
				Ver:  1,
				Type: auth.Type_Password,
				User: u,
				Pass: p,
			}); err != nil {
				return
			}
			r := new(auth.Result)
			if err = pbqp.Read(conn, r); err != nil {
				return
			}
			if r.Pass {
				fmt.Println(r.Welcome)
				break
			}
			fmt.Printf("\n%s\n\n", MsgAuthFail)
		}
	} else {
		var cipher []byte
		if cipher, err = vos.RsaEncrypt(pub, payload); err != nil {
			return
		}

		if err = pbqp.Write(conn, &auth.Auth{
			Ver:    1,
			Type:   auth.Type_PublicKey,
			Cipher: cipher,
		}); err != nil {
			return
		}
		r := new(auth.Result)
		if err = pbqp.Read(conn, r); err != nil {
			return
		}
		if err = vos.RsaVerify(pub, sig, r.Sig); err != nil {
			return
		}
		if !r.Pass {
			return ErrAuthFail
		}
		fmt.Println(r.Welcome)
	}
	return nil
}
