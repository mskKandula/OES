package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	stringsArr := strings.Split(currDir, string(os.PathSeparator))
	dirToServe := filepath.Join(stringsArr...)
	http.Handle("/", http.FileServer(http.Dir(dirToServe)))
	log.Fatal(http.ListenAndServe(":8887", nil))
}
