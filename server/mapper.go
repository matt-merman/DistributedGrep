package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type API int

type Couple struct {
	Key   string
	Value int
}

func (a *API) Mapper(sentence string, reply *[1000]Couple) error {

	v := strings.Split(sentence, " ")
	var couple Couple
	len := len(v)
	for i := 0; i < len; i++ {

		couple.Key = v[i]
		couple.Value = 1

		reply[i] = couple

	}

	return nil
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Distributed Grep Mapper 1.0\nUsage: go run mapper.go [port]\n")
		os.Exit(1)
	}

	port, _ := strconv.Atoi(os.Args[1])

	api := new(API)
	server := rpc.NewServer()
	err := server.RegisterName("API", api)
	if err != nil {
		log.Fatal("error registering API", err)
	}

	listener, err := net.Listen("tcp", ":"+os.Args[1])

	if err != nil {
		log.Fatal("Listener error", err)
	}

	log.Printf("\nMapper is listening on port %d", port)
	server.Accept(listener)
}
