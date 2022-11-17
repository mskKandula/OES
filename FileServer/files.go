package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("../OES/")))
	log.Fatal(http.ListenAndServe(":8887", nil))
}
