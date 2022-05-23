package main

import (
	"encoding/json"
	"github.com/google/uuid"
	ps "github.com/ppetar33/ars-project/poststore"
	"io"
	"net/http"
)

func decodeBody(r io.Reader) (*ps.Service, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt ps.Service
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeBodyConfig(r io.Reader) (*ps.Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt ps.Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId() string {
	return uuid.New().String()
}

func decodeConfigBody(r io.Reader) (*ps.Service, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var config *ps.Service
	if err := dec.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
