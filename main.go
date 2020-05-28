package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	go func() {
		fmt.Println("Service started!")
		fmt.Println(fmt.Sprintf("HTTP URL: http://%s:%d", host, port))

		srv.ListenAndServe()
	}()

	// Gracefully shutdown
	c := make(chan os.Signal, 1)

	// - catch quit via SIGINT (Ctrl+C) and SIGTERM (Kill command by the OS)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// - block until we receive our signal.
	oscall := <-c
	fmt.Println(fmt.Sprintf("Signal received:%+v", oscall))

	// - create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	// - doesn't block if no connections, but will otherwise wait until the timeout deadline.
	srv.Shutdown(ctx)
	fmt.Println("Shutting down service!")
	os.Exit(0)
}

// GetUserHandler ...
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user := map[string]interface{}{
		"id":   1,
		"name": "John Doe",
	}

	ub, _ := json.Marshal(user)

	time.Sleep(5 * time.Second)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ub)
}
