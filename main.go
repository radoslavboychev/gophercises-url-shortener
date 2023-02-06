package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	urlshort "github.com/radoslavboychev/gophercises/url-shortener/handlers"
)

var (
	pathsFile = flag.String("pathsFile", "paths.yml", "File containing shortened paths to URLs")
)

func getFileBytes(fileName string) []byte {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open %s", fileName)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		log.Fatalf("Could not open %s", fileName)
	}

	return buf.Bytes()

}

func main() {
	mux := defaultMux()

	flag.Parse()

	ext := filepath.Ext(*pathsFile)

	var handler http.Handler
	var err error
	if ext == ".yml" {
		handler, err = urlshort.YAMLHandler(getFileBytes(*pathsFile), mux)
		if err != nil {
			panic(err)
		}
	} else if ext == ".json" {
		handler, err = urlshort.JSONHandler(getFileBytes(*pathsFile), mux)
		if err != nil {
			panic(err)
		}
	} else {
		log.Fatal("Paths file needs to be either a YAML, a JSON or a bolt DB file")
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
