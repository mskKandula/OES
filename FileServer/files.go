package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	stringsArr := strings.Split(currDir, "/")
	dirToServe := "/" + stringsArr[1] + "/" + stringsArr[2] + "/"
	http.Handle("/", http.FileServer(http.Dir(dirToServe)))
	log.Fatal(http.ListenAndServe(":8887", nil))
}
