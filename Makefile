MAPPER = 2
PORT_SERVER = 54123
PORT_1 = 54203
PORT_2 = 54204


mapper:

	# go run server/mapper.go $(PORT_1) &
	# go run server/mapper.go $(PORT_2) &

reducer:

server:

	go run server/server.go 54123

client:
