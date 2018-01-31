package client

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/izhamoidsin/gedis/storage"
)

func handleGetResult(response *http.Response, err error) (storage.Storable, bool, error) {
	if err != nil {
		return nil, false, err
	}
	if response.StatusCode == http.StatusOK {
		val, err := parseStorableFormResponseBody(response)
		return val, true, err
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}

	return nil, false, errors.New("Unexpected response status code" + strconv.Itoa(response.StatusCode))
}

func parseStorableFormResponseBody(r *http.Response) (resp storage.Storable, err error) {
	var luckyString string
	var luckyArray []string
	var luckyDict map[string]string

	// TODO add parsing of TTL based header, e.g. Expire-In or Expire-At

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err = r.Body.Close(); err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &luckyString); err == nil {
		return luckyString, nil
	}

	if err = json.Unmarshal(body, &luckyArray); err == nil {
		return luckyArray, nil
	}

	if err = json.Unmarshal(body, &luckyDict); err == nil {
		return luckyDict, nil
	}

	return nil, err
}
