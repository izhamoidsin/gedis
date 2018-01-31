package client

import (
	"log"
	"testing"
	"time"

	"github.com/izhamoidsin/gedis/server"
	"github.com/izhamoidsin/gedis/storage"
)

var storageRegistry storage.Storage = storage.InitSyncStorage(time.Minute)
var srv = server.CreateServer(storageRegistry)

var client = CreateClient("localhost", 8088)

func TestRunTestServer(t *testing.T) {
	go func() { log.Fatal(srv.StartSerever(8088)) }()
}

func TestSaveGetUpdateAndDelete(t *testing.T) {
	key, value := "key", "value"
	newValue := "updated value"
	if error := client.AppendItem(key, value); error != nil {
		t.Error("Can not save item. " + error.Error())
	}
	if storedVal, exists, err := client.GetItem(key); err != nil || !exists || storedVal != value {
		t.Error("Can not get test value back")
	}
	if error := client.AppendItem(key, value); error == nil {
		t.Error("Save method allow to override resource")
	}
	if error := client.UpdateItem(key, newValue); error != nil {
		t.Error("Can not update test entry" + error.Error())
	}
	if storedVal, exists, err := client.GetItem(key); err != nil || !exists || storedVal != newValue {
		t.Error("Can not get updated test value back")
	}
	if error := client.DeleteItem(key); error != nil {
		t.Error("Can not delete item. " + error.Error())
	}
	if _, exists, err := client.GetItem(key); err != nil || exists {
		t.Error("Removed value is still returned" + err.Error())
	}
}

func TestNestedOps(t *testing.T) {
	arrayKey, array := "arr", []string{"Alpha", "Bravo", "Charlie"}
	dictKey, dict := "dic", map[string]string{
		"1": "One",
		"2": "Two",
	}

	if error := client.AppendItem(arrayKey, array); error != nil {
		t.Error("Can not save item. " + error.Error())
	}
	if error := client.AppendItem(dictKey, dict); error != nil {
		t.Error("Can not save item. " + error.Error())
	}
	if _, exists, error := client.GetItemByNestedIndex(arrayKey, "1"); !exists || error != nil {
		t.Error("Can not get element by index")
	}
	if _, exists, error := client.GetItemByNestedKey(dictKey, "2"); !exists || error != nil {
		t.Error("Can not get entry by key")
	}
}
