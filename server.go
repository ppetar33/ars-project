package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type postServer struct {
	data map[string][]*Config
}

func (ts *postServer) createConfigurationHandler(writer http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(writer, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	service, err := decodeBody(req.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	id := createId()
	service.Id = id
	service.Version = "1"
	ts.data[id] = service
	fmt.Println(service)

  renderJSON(writer, service)
}

func (ts *postServer) getConfigurationHandler(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, ok := ts.data[id]

	if !ok {
		err := errors.New("key not found")
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	if task.Version == version {
		fmt.Println(task)
		renderJSON(writer, task)
	} else {
		err := errors.New("key not found")
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

}

func (ts *postServer) getAllConfiugrationsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := make(map[string][]*Config)
	for k, v := range ts.data {
		allTasks[k] = v
	}
	renderJSON(w, allTasks)
}

func (ts *postServer) updateConfigurationHandler(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	contentType := req.Header.Get("Content-Type")
	mediatype, _, errParse := mime.ParseMediaType(contentType)

	if errParse != nil {
		http.Error(writer, errParse.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(writer, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	service, errDecode := decodeBody(req.Body)

	if errDecode != nil {
		http.Error(writer, errDecode.Error(), http.StatusBadRequest)
		return
	}

	ts.data[id] = append(ts.data[id], service...)

	renderJSON(writer, ts.data[id])
}

func (ts *postServer) delConfigurationHandler(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if v, ok := ts.data[id]; ok {
		delete(ts.data, id)
		renderJSON(writer, v)
	} else {
		err := errors.New("key not found")
		http.Error(writer, err.Error(), http.StatusNotFound)
	}
}
