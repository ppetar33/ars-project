package main

import (
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type postServer struct {
	data map[string]*Service
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
	ts.data[id] = service
	renderJSON(writer, service)
}

func (ts *postServer) getConfigurationHandler(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.data[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(writer, task)
}

func (ts *postServer) updateConfigurationHandler(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.data[id]
	contentType := req.Header.Get("Content-Type")
	mediatype, _, errParse := mime.ParseMediaType(contentType)

	if !ok {
		err := errors.New("key not found")
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

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

	task.Id = id
	service.Id = id

	renderJSON(writer, service)
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
