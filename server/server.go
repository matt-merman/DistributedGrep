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
const INITIAL_PORT = 54005
const PORT_REDUCER = 55001
const DEBUG = false

type Input struct {
	Word string
	File string
}

type API int

type Couple struct {
	Key   string
	Value int
}

type MapperInput struct {
	Text string
	Word string
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
	//split text in sentence (end with "\n")
	for scanner.Scan() {
		lines = append(lines, scanner.Text()+"\n")

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

func threadMapper(mapperPort int, sentence string, ch chan string, id int, word string) {

	//assign each part to a mapper with RPC
	var returnValue string

	mapper, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(mapperPort))
	defer mapper.Close()

	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	input := MapperInput{Text: sentence, Word: word}

	err = mapper.Call("API.Mapper", input, &returnValue)
	if err != nil {
		log.Fatal("Error in API.Mapper: ", err)
	}

	ch <- returnValue
}

func threadReducer(mapperPort int, input string, ch chan string) {

	//assign each part to a mapper with RPC
	var returnValue string

	mapper, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(PORT_REDUCER))
	defer mapper.Close()

	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	err = mapper.Call("API.Reducer", input, &returnValue)
	if err != nil {
		log.Fatal("Error in API.Reducer: ", err)
	}

	ch <- returnValue

}
func (a *API) Grep(input Input, reply *string) error {

	arrays, _ := openAndSplit(input.File)

	chMapper := make(chan string)
	defer close(chMapper)

	port := INITIAL_PORT
	for i := 0; i < NUMBER_MAPPER; i++ {

		go threadMapper(port, arrays[i], chMapper, i, input.Word)
		port += 1

	}

	var v string
	var ok bool
	var v1 string

	//listen on thread mapper channel until all mapper have terminated
	for i := 0; i < NUMBER_MAPPER; i++ {

		v, ok = <-chMapper
		if ok == false {
			break
		}
		v1 += v
		if DEBUG {
			fmt.Printf("Main Thread (%d) has received '%s' in\n%s\n", os.Getpid(), input.Word, v)
		}
	}

	//shuffle and sort phase
	//(empty)

	//reducer with RPC
	go threadReducer(port, v1, chMapper)
	vt, err := <-chMapper
	if !err {
		log.Fatal("Error in API.Reducer: ", err)
	}

	*reply = vt
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
