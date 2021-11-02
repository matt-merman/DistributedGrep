This is a particular implementation of the _grep_ from command line using _RPC_ and written in _Go_. _API.Grep_ returns the number of occurences for a word in a text. To execute both client and server:

```
go run server.go
go run client.go che test.txt
```

To verify the obtained result it is sufficient to run the command _grep -o 'che' test.txt | wc -l_ and compare the results. _API.grep_ will also count words like _checosa_ or _albecheoro_.

