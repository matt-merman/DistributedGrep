package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Input struct {
	Text string
	Word string
}

func main() {

	fmt.Printf("Distributed Grep Client\n")

	if len(os.Args) != 4 {
		fmt.Printf("Usage: go run client.go [word] [file] [port server]\n")
		os.Exit(1)
	}

	file := os.Args[2]
	word := os.Args[1]

	client, err := rpc.Dial("tcp", "localhost:"+os.Args[3])

	defer client.Close()

	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	var returnValue string
	err = client.Call("API.Grep", Input{file, word}, &returnValue)
	if err != nil {
		log.Fatal("Error in API.Grep: ", err)
	}

	fmt.Printf("Senteces that contain '%s' are:\n\n%s", word, returnValue)
}
