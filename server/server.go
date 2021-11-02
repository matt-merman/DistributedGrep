package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type Input struct {
	Word string
	File string
}

type API int

type couple struct {
	key   string
	value int
}

func mapper(id int, sentence string, word string, ch chan couple) {

	//strings.Count() counts even occurences in combined word as
	//strings.Count("quantitutti", "tutti") returns 1
	count := strings.Count(sentence, word)

	return_value := couple{key: word, value: count}
	ch <- return_value

}

func reducer(node map[int]couple, ch chan map[int]couple) {

	ch <- node

}

func openAndSplit(file string) (*bufio.Scanner, *os.File) {

	//open file
	f, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	//split file into senteces
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	return scanner, f

}

func (a *API) Grep(input Input, reply *int) error {

	sentences, fileOpened := openAndSplit(input.File)
	if fileOpened == nil {
		return nil
	}

	defer fileOpened.Close()

	//define a channel to exchange key:value with mapper
	chMapper := make(chan couple)
	defer close(chMapper)

	//define a channel to exchange key:value with reducer
	chReducer := make(chan map[int]couple)
	defer close(chReducer)

	//run a goroutine for each sentence
	var counterSentence int
	for sentences.Scan() {

		sentence := sentences.Text()
		go mapper(counterSentence, sentence, input.Word, chMapper)
		counterSentence++
	}

	if err := sentences.Err(); err != nil {
		fmt.Println(err)
	}

	//map where save all couple returned from mappers
	//mpCoupleMapper = [{0, ["word":1]}, {1, ["word":4]}, ...]
	mpCoupleMapper := make(map[int]couple)

	//listen on mapper channel until all mapper have terminated
	for i := 0; i < counterSentence; i++ {

		v, ok := <-chMapper
		if ok == false {
			break
		}

		mpCoupleMapper[i] = couple{v.key, v.value}
	}

	//shuffle and sort phase
	//(empty)

	//in our case there's a unique reducer (only one key)
	go reducer(mpCoupleMapper, chReducer)

	//map where save all couple returned from the unique reducer
	//mpCoupleReducer = [{0, ["word":1]}, {1, ["word":4]}, ...]
	mpCoupleReducer := make(map[int]couple)

	mpCoupleReducer, err := <-chReducer
	if err == false {
		return nil
	}

	var occurences int
	//main thread counts occurences
	for _, value := range mpCoupleReducer {
		occurences += value.value
	}

	*reply = occurences
	return nil
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Distributed Grep Server 1.0\nUsage: go run server.go [port]\n")
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
	log.Printf("\nServer is listening on port %d", port)
	server.Accept(listener)
}
