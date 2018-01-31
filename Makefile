all: dependencies clean test build

dependencies:
	go get -u github.com/gorilla/mux
	go get -u github.com/golang/sync/syncmap
	go get -u github.com/izhamoidsin/gedis/storage
	go get -u github.com/izhamoidsin/gedis/server

test:
	go test -cover ./...

clean:
	go clean ./...

build:
	go build  -v ./...
