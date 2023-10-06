# Golang REST API

## Description
This API reads data from the /test/input file and starts an X number of workers provided by 
a POST request. Each worker reads a number of lines from the input file, based on the total
number of workers and the number of lines in the file. The workers then process the data 
and write a log of their result to the database. The database is a MySQL database running in a
Docker container.

## Run instructions
1. run the Dockerfile inside the /build folder to create the database
2. run `go get -d ./...` to download all dependencies
3. run `go run cmd/firstProject/main.go` to start the server

## Endpoints
- **POST /api/v1** - launch new job
    - body: `"{numberOfWorkers": <integer greater than 0>}`
