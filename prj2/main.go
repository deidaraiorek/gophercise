package main

import (
	"fmt"
	"net/http"
	"os"
	"prj2/urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// yaml, err := os.ReadFile("./urls.yaml")
	// if err != nil {
	// 	panic(err)
	// }
	// yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	// if err != nil {
	// 	panic(err)
	// }

	jsonFile, err := os.ReadFile("./urls.json")
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler(jsonFile, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8081")
	http.ListenAndServe(":8081", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
