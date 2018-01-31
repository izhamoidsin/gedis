package server

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/izhamoidsin/gedis/storage"
)

func parseJSONFormRequestBody(r *http.Request) (storage.Storable, error) {
	var luckyString string
	var luckyArray []string
	var luckyDict map[string]string

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &luckyString); err == nil {
		return luckyString, nil
	}

	if err := json.Unmarshal(body, &luckyArray); err == nil {
		return luckyArray, nil
	}

	if err := json.Unmarshal(body, &luckyDict); err == nil {
		return luckyDict, nil
	}

	return nil, errors.New("Entity is unprocessable")
}

func getPathVars(r *http.Request) (key string, subKey string, index int) {
	vars := mux.Vars(r)
	key = vars["key"]
	subKey = vars["subKey"]
	if indexStr, exists := vars["index"]; exists {
		if indexNumber, err := strconv.Atoi(indexStr); err == nil {
			index = indexNumber
		}
	}

	return key, subKey, index
}

func respondWithJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func respondWithExpireIn(w http.ResponseWriter) {
	// TODO add TTL based header, e.g. Expire-At or Expire-In
	w.Header().Set("Expire-At", "")
}
