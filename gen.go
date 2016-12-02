package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

func generateInit() {
	err := os.MkdirAll("./grifts", 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(path.Join("./grifts", "example.go"), []byte(initTmpl), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
