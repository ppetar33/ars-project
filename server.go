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

	res := map[string]*[]Config{}

	for _, m := range task.Data {
		for k, v := range service.Data {
			res[k] = append(res[k], v)
		}
	}

	renderJSON(writer, res)
}

func merge(ms ...map[string]*[]Config) map[string]*[]Config {
	res := map[string]*[]Config{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = append(res[k], v)
		}
	}
	return res
}
