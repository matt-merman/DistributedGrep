#!/bin/sh
# bash mapperCreation.sh maxNumberMapper portBase
#portBase=5000
#maxNumberMapper=10

#$1 = 10
#$2 = 100
for ((i = 0 ; i < 10 ; i++)); do
  newPort=$((100 + $i))
  go run ./server/mapper.go $newPort
done
