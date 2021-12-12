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

const DEBUG = false
const NUMBER_MAPPER = 2

var portMapper int
var portReducer int

type API int

type Input struct {
	Text string
	Word string
}

func openAndSplit(name string) [NUMBER_MAPPER]string {

	var sentences [NUMBER_MAPPER]string

	file, err := os.Open(name)
	defer file.Close()

	if err != nil {
		log.Fatal("Error in openAndSplit: ", err)
	}

	//----------------

	var lines []string
	scanner := bufio.NewScanner(file)
	//split text in sentence (end with "\n")
	for scanner.Scan() {
		lines = append(lines, scanner.Text()+"\n")

	}

	//divide array into NUMBER_MAPPER parts
	numberWords := len(lines)
	dimensionPart := numberWords / NUMBER_MAPPER
	initialPart := dimensionPart

	k := 0
	for i := 0; i < numberWords; i++ {

		if i > initialPart {

			initialPart += dimensionPart
			k++
		}
		sentences[k] += lines[i]
	}

	return sentences
}

func threadMapper(port int, sentence string, channel chan string, word string) {

	//assign each part to a mapper with RPC
	var returnValue string

	mapper, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(port))
	defer mapper.Close()

	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	input := Input{Text: sentence, Word: word}

	err = mapper.Call("API.Mapper", input, &returnValue)
	if err != nil {
		log.Fatal("Error in API.Mapper: ", err)
	}

	channel <- returnValue
}

func threadReducer(port int, input string, channel chan string) {

	//assign each part to a mapper with RPC
	var returnValue string

	mapper, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(port))
	defer mapper.Close()

	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	err = mapper.Call("API.Reducer", input, &returnValue)
	if err != nil {
		log.Fatal("Error in API.Reducer: ", err)
	}

	channel <- returnValue

}

func (a *API) Grep(input Input, reply *string) error {

	//1. SPLIT INPUT
	arrayText := openAndSplit(input.Text)

	//2. MAPPER
	channel := make(chan string)
	defer close(channel)

	for index := 0; index < NUMBER_MAPPER; index++ {

		go threadMapper(portMapper, arrayText[index], channel, input.Word)
		portMapper += 1

	}

	var sentencesMapper string

	//listen on thread mapper channel until all mapper have terminated
	for i := 0; i < NUMBER_MAPPER; i++ {

		sentence, err := <-channel
		if !err {
			log.Fatal("Error in API.Grep: ", err)
		}
		sentencesMapper += sentence

		if DEBUG {
			fmt.Printf("Main Thread (%d) has received '%s' in\n%s\n", os.Getpid(), input.Word, sentence)
		}
	}

	//3. SHUFFLE & SORT
	//(empty)

	//4. REDUCER
	go threadReducer(portReducer, sentencesMapper, channel)
	sentencesReducer, err := <-channel
	if !err {
		log.Fatal("Error in API.Grep: ", err)
	}

	*reply = sentencesReducer
	return nil

}

func main() {

	if len(os.Args) != 4 {
		log.Fatal("Usage: go run server.go [port server] [port mapper] [port reducer]\n")
	}

	port, _ := strconv.Atoi(os.Args[1])
	portMapper, _ = strconv.Atoi(os.Args[2])
	portReducer, _ = strconv.Atoi(os.Args[3])

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

	log.Printf("\nDistributed Grep Server is listening on port %d", port)

	server.Accept(listener)

}
