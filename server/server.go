package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

const NUMBER_MAPPER = 2
const MAX_WORD = 1000
const INITIAL_PORT = 54322
const DEBUG = true

type Input struct {
	Word string
	File string
}

type API int

type Couple struct {
	Key   string
	Value int
}

func openAndSplit(file string) ([NUMBER_MAPPER]string, *os.File) {

	var sentences [NUMBER_MAPPER]string

	//open file
	f, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
		return sentences, nil
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	//divide array into NUMBER_MAPPER parts
	number_words := len(lines)

	dimensionPart := number_words / NUMBER_MAPPER
	initialPart := dimensionPart

	k := 0
	for i := 0; i < number_words; i++ {

		if i > initialPart {

			initialPart += dimensionPart
			k++
		}
		sentences[k] += lines[i] + " "
	}

	return sentences, f

}

func threadMapper(mapperPort int, sentence string, ch chan Couple, id int) {

	//assign each part to a mapper with RPC
	var returnValue [MAX_WORD]Couple

	mapper, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(mapperPort))
	defer mapper.Close()

	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	err = mapper.Call("API.Mapper", sentence, &returnValue)
	if err != nil {
		log.Fatal("Error in API.Mapper: ", err)
	}

	fmt.Printf("Thread %d Running\n", id)

	for i := 0; i < MAX_WORD; i++ {

		if returnValue[i].Value != 1 {
			break
		}
		if DEBUG {
			fmt.Printf("(%d, %s, %d)\n", id, returnValue[i].Key, returnValue[i].Value)
		}
	}

	for i := 0; i < MAX_WORD; i++ {

		ch <- returnValue[i]
	}
}

func (a *API) Grep(input Input, reply *string) error {

	arrays, _ := openAndSplit(input.File)

	chMapper := make(chan Couple)
	defer close(chMapper)

	port := INITIAL_PORT
	for i := 0; i < NUMBER_MAPPER; i++ {

		go threadMapper(port, arrays[i], chMapper, i)
		port += 1

	}

	//listen on thread mapper channel until all mapper have terminated
	for i := 0; i < NUMBER_MAPPER*MAX_WORD; i++ {

		v, ok := <-chMapper
		if ok == false {
			break
		}

		if DEBUG && v.Value == 1 {
			fmt.Printf("Main Thread: (%d, %s, %d)\n", os.Getpid(), v.Key, v.Value)
		}
	}

	//shuffle and sort phase
	//(empty)

	//reducer with RPC
	//...

	*reply = "ciao"
	return nil

}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Distributed Grep Server 1.0\nUsage: go run server.go [port]\n")
		os.Exit(1)
	}

	port, _ := strconv.Atoi(os.Args[1]) //54123

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

	log.Printf("\nServer is listening on port %d", port)
	server.Accept(listener)

}
