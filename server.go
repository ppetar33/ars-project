package main

import (
	"errors"
	"github.com/gorilla/mux"
	ps "github.com/ppetar33/ars-project/poststore"
	"mime"
	"net/http"
)

type postServer struct {
	store *ps.PostStore
}

func (ts *postServer) createConfigurationHandler(writer http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	idempotencyId := req.Header.Get("x-idempotency-key")

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

	existService, err := ts.store.FindConfByIdempotency(idempotencyId)
	if existService != nil {
		post, err := ts.store.Post(service, idempotencyId)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		renderJSON(writer, post)
	} else {
		renderJSON(writer, service)
	}

}

func (ts *postServer) updateConfigurationHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	id := mux.Vars(req)["id"]
	idempotencyId := req.Header.Get("x-idempotency-key")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeConfigBody(req.Body)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	existService, err := ts.store.FindConfByIdempotency(idempotencyId)
	if existService != nil {
		rt.Id = id
		service, err := ts.store.Update(rt, idempotencyId)

		if err != nil {
			http.Error(w, "Given config version already exists! ", http.StatusBadRequest)
			return
		}

		w.Write([]byte(service.Id))
	} else {
		renderJSON(w, rt)
	}

}

func (ts *postServer) getAllConfigurationsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := ts.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
}

func (ts *postServer) findConfigurationsByLabels(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	ver := mux.Vars(req)["version"]

	configs, err := decodeBodyConfig(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, ok := ts.store.FindByLabels(id, ver, configs)

	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if task == nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *postServer) getConfigurationByIdAndVersion(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	ver := mux.Vars(req)["version"]

	task, ok := ts.store.Get(id, ver)

	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if task == nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *postServer) delConfigurationHandler(writer http.ResponseWriter, req *http.Request) {
	ver := mux.Vars(req)["version"]
	id := mux.Vars(req)["id"]
	task, ok := ts.store.Delete(id, ver)
	if ok != nil {
		err := errors.New("not found")
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(writer, task)
}

func (ts *postServer) getConfigurationById(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.store.FindConfVersions(id)
	if ok != nil {
		err := errors.New("not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if task == nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}
