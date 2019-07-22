package main

import (
	"io/ioutil"
	"net"
	"os"

	"github.com/ddosakura/vos"
	"github.com/ddosakura/vos/proto/auth"
)

func main() {
	r, _ := os.Open(os.Args[2])
	bs, _ := ioutil.ReadAll(r)
	priv, _ := vos.RsaPriv(bs)
	o := vos.New()
	o.Auth = func(req *auth.Auth, res *auth.Result) (string, error) {
		switch req.Type {
		case auth.Type_Password:
			if req.User == "user" && req.Pass == "123456" {
				res.Pass = true
				res.Welcome = `Success`
				return req.User, nil
			}
		case auth.Type_PublicKey:
			var err error
			if res.Sig, err = vos.RsaSign(priv, []byte("root")); err != nil {
				return "", err
			}
			if bs, err := vos.RsaDecrypt(priv, req.Cipher); err != nil {
				return "", err
			} else if string(bs) == "vos" {
				res.Pass = true
				res.Welcome = `Success`
				if err != nil {
					return "", err
				}
				return "root", nil
			}
		default:
			return "", vos.ErrAuthProtocol
		}
		return "", nil
	}
	l, err := net.Listen("unix", os.Args[1])
	if err != nil {
		o.Error(err)
	}
	o.Run(l)
}
