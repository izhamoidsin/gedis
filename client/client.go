package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/izhamoidsin/gedis/storage"
)

// TODO verify all errors supressing with _
// TODO refactor error handling

// GedisClient is go lang client to Gedis Server. Wraps HTTP calls and provide
// a native API
type GedisClient struct {
	host    string
	port    int
	strPort string
}

// CreateClient call creates a new instance of client and initializes it
func CreateClient(host string, port int) *GedisClient {
	client := new(GedisClient)
	client.host = host
	client.port = port
	client.strPort = strconv.Itoa(port)
	return client
}

func (client *GedisClient) fullURL(path string) string {
	return "http://" + client.host + ":" + client.strPort + "/" + path
}

// GetKeys call retruns slice of all the keys stored in Gedis at the moment
// or an error if appeared
func (client *GedisClient) GetKeys() ([]string, error) {
	respose, err := http.Get(client.fullURL("keys"))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(io.LimitReader(respose.Body, 1048576))
	if err != nil {
		return nil, err
	}

	if err = respose.Body.Close(); err != nil {
		return nil, err
	}

	var keys []string
	if err = json.Unmarshal(body, &keys); err == nil {
		return keys, nil
	}
	return nil, err
}

// GetItem ...
func (client *GedisClient) GetItem(key string) (storage.Storable, bool, error) {
	response, err := http.Get(client.fullURL("entries/" + key))
	return handleGetResult(response, err)
}

// UpdateItem ...
func (client *GedisClient) UpdateItem(key string, item storage.Storable) error {
	bts, err := json.Marshal(item)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		http.MethodPut,
		client.fullURL("entries/"+key),
		bytes.NewBuffer(bts),
	)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if response.StatusCode != http.StatusNoContent {
		err = errors.New("Unexpected response status code")
	}
	return err
}

// AppendItem ...
func (client *GedisClient) AppendItem(key string, item storage.Storable) error {
	bts, _ := json.Marshal(item)
	response, err := http.Post(client.fullURL("entries/"+key), "application/json; charset=UTF-8", bytes.NewBuffer(bts))

	if response.StatusCode != http.StatusCreated {
		err = errors.New("Unexpected response status code")
	}

	return err
}

// DeleteItem ...
func (client *GedisClient) DeleteItem(key string) error {
	request, _ := http.NewRequest(
		http.MethodDelete,
		client.fullURL("entries/"+key),
		bytes.NewBuffer([]byte{}),
	)
	response, error := http.DefaultClient.Do(request)

	if response.StatusCode != http.StatusNoContent {
		error = errors.New("Unexpected response status code")
	}

	return error
}

// GetItemByNestedIndex ...
func (client *GedisClient) GetItemByNestedIndex(key string, index string) (storage.Storable, bool, error) {
	response, err := http.Get(client.fullURL("entries/" + key + "/elements/" + index)) // TODO make index numeric
	return handleGetResult(response, err)
}

// GetItemByNestedKey ...
func (client *GedisClient) GetItemByNestedKey(key string, subKey string) (storage.Storable, bool, error) {
	response, err := http.Get(client.fullURL("entries/" + key + "/entries/" + subKey))
	return handleGetResult(response, err)
}
