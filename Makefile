SERVER = 59000
MAPPER1 = 59020
MAPPER2 = 59021
REDUCER = 59003

server_run:

	cd ./server; go build -o mapper.out mapper.go
	cd ./server; ./mapper.out $(MAPPER1) &
	cd ./server; ./mapper.out $(MAPPER2) &
	
	cd ./server; go build -o reducer.out reducer.go
	cd ./server; ./reducer.out $(REDUCER) &

	cd ./server; go build -o server.out server.go
	cd ./server; ./server.out $(SERVER) $(MAPPER1) $(REDUCER) &

client_run:

	go build -o client.out client.go
	./client.out cat catas.txt $(SERVER)

kill:

	pkill server.out || true
	pkill mapper.out || true
	pkill reducer.out || true

clean:

	go clean
	cd ./server; rm mapper.out; rm reducer.out; rm server.out;
	rm client.out

