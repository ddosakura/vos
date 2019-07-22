package main

import (
	"crypto/rsa"
	"io/ioutil"
	"net"
	"os"

	"github.com/chzyer/readline"
	"github.com/ddosakura/vos"
	"github.com/ddosakura/vos/vterm"
)

func main() {
	conn, err := net.Dial("unix", os.Args[1])
	if err != nil {
		panic(err)
	}
	var pub *rsa.PublicKey
	if len(os.Args) > 2 {
		r, _ := os.Open(os.Args[2])
		bs, _ := ioutil.ReadAll(r)
		pub, _ = vos.RsaPub(bs)
	}

	if err := vterm.Connect(conn, pub, []byte("vos"), []byte("root")); err != nil {
		panic(err)
	}

	// TODO: SMUX 多路复用
	cli, err := readline.NewRemoteCli(conn)
	if err != nil {
		panic(err)
	}
	if err = cli.Serve(); err != nil {
		panic(err)
	}
}
