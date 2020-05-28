package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	log.Println("Service started!")
	log.Println(fmt.Sprintf("HTTP URL: http://%s:%d", host, port))

	srv.ListenAndServe()
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
