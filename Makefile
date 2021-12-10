SERVER = 57123
MAPPER_1 = 56008
MAPPER_2 = 56009
MAPPER_3 = 56010
REDUCER = 57001

server_run:

	cd ./server; go build -o mapper.out mapper.go
	cd ./server; ./mapper.out $(MAPPER_1) > /dev/null 2>&1 &
	cd ./server; ./mapper.out $(MAPPER_2) > /dev/null 2>&1 &
	
	cd ./server; go build -o reducer.out reducer.go
	cd ./server; ./reducer.out $(REDUCER) > /dev/null 2>&1 &

	cd ./server; go build -o server.out server.go
	cd ./server; ./server.out $(SERVER) > /dev/null 2>&1 &

client_run:

	go build -o client.out client.go
	./client.out che test.txt $(SERVER)

kill:

clean:

	go clean
	cd ./server; rm mapper.out; rm reducer.out; rm server.out;
	rm client.out

