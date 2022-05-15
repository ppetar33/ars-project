package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type postServer struct {
	data map[string][]*Config `json:"data"`
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

func (ts *postServer) getAllConfiugrationsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*Config{}
	for _, v := range ts.data {
		allTasks = append(allTasks, v)
	}
	renderJSON(w, allTasks)
}

func (ts *postServer) updateConfigurationHandler(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.data[id]
	//contentType := req.Header.Get("Content-Type")
	//mediatype, _, errParse := mime.ParseMediaType(contentType)

	if !ok {
		err := errors.New("key not found")
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	//
	//if errParse != nil {
	//	http.Error(writer, errParse.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//if mediatype != "application/json" {
	//	err := errors.New("expect application/json Content-Type")
	//	http.Error(writer, err.Error(), http.StatusUnsupportedMediaType)
	//	return
	//}
	//
	service, errDecode := decodeBody(req.Body)
	//
	if errDecode != nil {
		http.Error(writer, errDecode.Error(), http.StatusBadRequest)
		return
	}
	//
	//task.Id = id
	//service.Id = id

	fmt.Println(task.Entries)    // staro
	fmt.Println(service.Entries) // novo
	//fmt.Println(service.Entries)

	for k, v := range service.Entries {
		task.Entries[k] = v
	}

	fmt.Println(task.Entries)

	renderJSON(writer, task)
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
