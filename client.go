package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Input struct {
	Word string
	File string
}

func main() {

	if len(os.Args) != 4 {
		fmt.Printf("Distributed Grep Client 1.0\nUsage: go run client.go [word] [file] [port]\n")
		os.Exit(1)
	}

	file := os.Args[2]
	word := os.Args[1]

	var returnValue int

	client, err := rpc.Dial("tcp", "localhost:"+os.Args[3])

	defer client.Close()

	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	err = client.Call("API.Grep", Input{word, file}, &returnValue)
	if err != nil {
		log.Fatal("Error in API.Grep: ", err)
	}

	fmt.Printf("Number of occurences for '%s' in '%s' is %d\n", word, file, returnValue)
}
