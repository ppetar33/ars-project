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
		data: map[string]*Config{},
	}

	router.HandleFunc("/conf/create/", server.createConfigurationHandler).Methods("POST")
	router.HandleFunc("/conf/{id}/", server.getConfigurationHandler).Methods("GET")
	router.HandleFunc("/conf/", server.getAllConfiugrationsHandler).Methods("GET")
	router.HandleFunc("/conf/extend/{id}/", server.updateConfigurationHandler).Methods("PUT")
	router.HandleFunc("/conf/delete/{id}/", server.delConfigurationHandler).Methods("DELETE")

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
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
