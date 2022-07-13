package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am service b\n"))
	})

	http.ListenAndServe(":8500", nil)
}
