package main

import (
	"log"

	"github.com/izhamoidsin/gedis/server"
	"github.com/izhamoidsin/gedis/storage"
)

// Runs the server according to the config
func main() {
	var storage storage.Storage = storage.InitSyncStorage(ttl)
	var server = server.CreateServer(storage)

	log.Fatal(server.StartSerever(port))
}
