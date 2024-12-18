package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Rule Engine API Service is running!")
	})
	http.ListenAndServe(":8003", nil)
}
