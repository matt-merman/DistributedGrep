package main

type couple struct {
	key   string
	value int
}

func reducer(node map[int]couple, ch chan map[int]couple) {

	ch <- node

}
