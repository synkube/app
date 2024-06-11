package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Start a goroutine to print a message every 5 seconds
	go func() {
		for {
			fmt.Println("Message printed every 5 seconds")
			fmt.Println("Message printed every 10 seconds")
			time.Sleep(5 * time.Second)
		}
	}()

	// Handle root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've reached the Go web server!")
	})

	// Start the web server
	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
