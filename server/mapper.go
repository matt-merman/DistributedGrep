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

type Input struct {
	Text string
	Word string
}

const DEBUG = false

func (a *API) Mapper(input Input, reply *string) error {

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
		log.Fatal("Usage: go run mapper.go [port]\n")
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

	log.Printf("\nDistributed Grep Mapper is listening on port %d", port)
	server.Accept(listener)
}
