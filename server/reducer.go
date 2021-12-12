package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

type API int

const DEBUG = false

func (a *API) Reducer(input string, reply *string) error {

	if DEBUG {
		fmt.Printf("Reducer (%d) has received \n%s\n", os.Getpid(), input)
	}

	*reply = input
	return nil
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Usage: go run reducer.go [port]\n")
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

	log.Printf("\nDistributed Grep Reducer is listening on port %d", port)
	server.Accept(listener)
}
