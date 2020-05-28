package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	host = "0.0.0.0"
	port = 8080
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/user", GetUserHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: mux,
	}

	// Run http server in go routine
	go func() {
		log.Println("Service started!")
		log.Println(fmt.Sprintf("HTTP URL: http://%s:%d", host, port))

		srv.ListenAndServe()
	}()

	ImplementGraceful(srv, time.Duration(60)*time.Second)
}

// ImplementGraceful will implement graceful shutdown to the given
// http server object. We will wait for SIGINT & SIGTERM signal
// before initiating the graceful shutdown sequence.
//
// If the existing connection no yet finished after the given timeout,
// we will forcefully shutdown the server.
func ImplementGraceful(s *http.Server, timout time.Duration) {
	// Make channel, and listen for SIGINT & SIGTERM
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until receive the signal
	oscall := <-c
	log.Println(fmt.Sprintf("Signal received:%+v", oscall))

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), timout)
	defer cancel()

	// Gracefully shutdown the http server
	s.Shutdown(ctx)

	// Exiting the program
	log.Println("Shutting down service!")
	os.Exit(0)
}

// GetUserHandler implementation
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user := map[string]interface{}{
		"id":   1,
		"name": "John Doe",
	}

	ub, _ := json.Marshal(user)

	// Add 5 second timeout to simulate slow upstream services
	time.Sleep(5 * time.Second)

	w.Header().Set("Content-Type", "application/json")
	w.Write(ub)
}
