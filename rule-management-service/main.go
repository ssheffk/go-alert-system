package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Rule Management API Service is running!")
	})
	http.ListenAndServe(":8004", nil)
}
