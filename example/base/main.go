package main

import (
	"net"
	"os"

	"github.com/ddosakura/vos"
)

func main() {
	o := vos.New()
	l, err := net.Listen("unix", os.Args[1])
	if err != nil {
		o.Error(err)
	}
	o.Run(l)
}
