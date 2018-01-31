package server

import (
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/izhamoidsin/gedis/storage"
)

// GedisServer ...
type GedisServer struct {
	startTime time.Time
	storage   storage.Storage
}

// CreateServer ...
func CreateServer(storage storage.Storage) *GedisServer {
	server := new(GedisServer)
	server.storage = storage
	return server
}

// StartSerever ...
func (server *GedisServer) StartSerever(port int) error {
	server.startTime = time.Now()

	// used gorilla mux router because it reduces boilerplate code of http methods & paths matching
	router := mux.NewRouter()
	router.HandleFunc("/heartbeat", server.heartbeat).Methods(http.MethodGet, http.MethodHead)
	router.HandleFunc("/keys", server.keys).Methods(http.MethodGet)
	router.HandleFunc("/entries/{key}", server.getItem).Methods(http.MethodGet)
	router.HandleFunc("/entries/{key}", server.chechItemPresense).Methods(http.MethodHead)
	router.HandleFunc("/entries/{key}", server.putItem).Methods(http.MethodPut)
	router.HandleFunc("/entries/{key}", server.appendItem).Methods(http.MethodPost)
	router.HandleFunc("/entries/{key}", server.deleteItem).Methods(http.MethodDelete)
	router.HandleFunc("/entries/{key}/elements/{ind}", server.getByNestedIndex).Methods(http.MethodGet)
	router.HandleFunc("/entries/{key}/entries/{subKey}", server.getByNestedKey).Methods(http.MethodGet)

	http.Handle("/", router)
	log.Println("Starting server @ port " + strconv.Itoa(port))
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (server *GedisServer) heartbeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I'm ok sinse "+server.startTime.Format(time.RFC850))
	respondWithJSON(w)
}

func (server *GedisServer) keys(w http.ResponseWriter, r *http.Request) {
	keys := server.storage.GetAllKeys()
	json.NewEncoder(w).Encode(keys)
	respondWithJSON(w)
}

func (server *GedisServer) putItem(w http.ResponseWriter, r *http.Request) {
	key, _, _ := getPathVars(r)
	if newValue, err := parseJSONFormRequestBody(r); err == nil {
		if operationForbidden := server.storage.UpdateValueByKey(key, newValue); operationForbidden == nil {
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, operationForbidden.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (server *GedisServer) appendItem(w http.ResponseWriter, r *http.Request) {
	key, _, _ := getPathVars(r)
	if newValue, err := parseJSONFormRequestBody(r); err == nil {
		if operationForbidden := server.storage.AppendNewValue(key, newValue); operationForbidden == nil {
			w.WriteHeader(http.StatusCreated)
			// TODO add Location header & make response compliant to rfc2616
		} else {
			http.Error(w, operationForbidden.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (server *GedisServer) deleteItem(w http.ResponseWriter, r *http.Request) {
	key, _, _ := getPathVars(r)
	server.storage.DeleteValueByKey(key) // TODO handle deleted flag
	w.WriteHeader(http.StatusNoContent)
}

func (server *GedisServer) chechItemPresense(w http.ResponseWriter, r *http.Request) {
	key, _, _ := getPathVars(r)
	if _, ok := server.storage.GetValueByKey(key); ok {
		return
	}
	http.NotFound(w, r)
}

func (server *GedisServer) getItem(w http.ResponseWriter, r *http.Request) {
	key, _, _ := getPathVars(r)
	if val, ok := server.storage.GetValueByKey(key); ok {
		json.NewEncoder(w).Encode(val.Entity)
		respondWithJSON(w)
	} else {
		http.NotFound(w, r)
	}
}

func (server *GedisServer) getByNestedKey(w http.ResponseWriter, r *http.Request) {
	key, subKey, _ := getPathVars(r)
	if val, exists, error := server.storage.GetNestedValueByKeyAndSubkey(key, subKey); error == nil && exists {
		json.NewEncoder(w).Encode(val.Entity)
		respondWithJSON(w)
	} else if !exists && error == nil {
		http.NotFound(w, r)
	} else { // TODO identify bad op error
		http.Error(w, error.Error(), http.StatusBadRequest)
	}
}

func (server *GedisServer) getByNestedIndex(w http.ResponseWriter, r *http.Request) {
	key, _, index := getPathVars(r)
	// TODO handle errors
	if val, exists, error := server.storage.GetNestedValueByKeyAndIndex(key, index); error == nil && exists {
		json.NewEncoder(w).Encode(val.Entity)
		respondWithJSON(w)
	} else if !exists && error == nil {
		http.NotFound(w, r)
	} else { // TODO identify bad op error
		http.Error(w, error.Error(), http.StatusBadRequest)
	}
}
