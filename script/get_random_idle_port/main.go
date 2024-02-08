package main

import (
	"fmt"
	"net"
	"strconv"
)

func main() {
	fmt.Print(GetRandomIdlePort())
}

func GetRandomIdlePort() string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	return strconv.Itoa(port)
}
