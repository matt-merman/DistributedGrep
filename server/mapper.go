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

type MapperInput struct {
	Text string
	Word string
}

const DEBUG = true

func (a *API) Mapper(input MapperInput, reply *string) error {

	v := strings.Split(input.Text, "\n")
	len := len(v)
	for i := 0; i < len; i++ {

		count := strings.Count(v[i], input.Word)
		if count != 0 {
			if DEBUG {
				fmt.Printf("Mapper (%d) has founded '%s' in \n%s\n", os.Getpid(), input.Word, v[i])
			}
			*reply += v[i] + "\n"
		}

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
