package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	server := postServer{
		data: map[string]*Service{},
	}

	router.HandleFunc("/conf/create/", count(server.createConfigurationHandler)).Methods("POST")
	router.HandleFunc("/conf/{id}/{version}/", count(server.getConfigurationHandler)).Methods("GET")
	router.HandleFunc("/conf/", count(server.getAllConfigurationsHandler)).Methods("GET")
	router.HandleFunc("/conf/extend/{id}/{version}/", count(server.updateConfigurationHandler)).Methods("PUT")
	router.HandleFunc("/conf/delete/{id}/{version}/", count(server.delConfigurationHandler)).Methods("DELETE")
	router.Path("/metrics").Handler(metricsHandler())

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}
