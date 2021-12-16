SERVER = 12345
MAPPER1 = 12346
MAPPER2 = 12347
MAPPER3 = 12348
REDUCER = 12349

server_run_client_build:

	cd ./code/server; go build -o mapper.out mapper.go
	cd ./code/server; ./mapper.out $(MAPPER1) &
	cd ./code/server; ./mapper.out $(MAPPER2) &
	cd ./code/server; ./mapper.out $(MAPPER3) &
	
	cd ./code/server; go build -o reducer.out reducer.go
	cd ./code/server; ./reducer.out $(REDUCER) &

	cd ./code/server; go build -o server.out server.go
	cd ./code/server; ./server.out $(SERVER) $(MAPPER1) $(REDUCER) &

	go build -o ./code/client.out ./code/client.go

client_run:

	./code/client.out cats cats.txt $(SERVER)

kill:

	pkill server.out  || true
	pkill mapper.out  || true
	pkill reducer.out || true

clean:

	go clean
	cd ./code/server; rm mapper.out; rm reducer.out; rm server.out;
	rm ./code/client.out

grep:

	cd ./code/server; cat cats.txt | grep cats



