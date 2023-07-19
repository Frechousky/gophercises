package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"frechousky/urlshort"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func parseCli() (*string, *string) {
	yamlFile := flag.String("yml", "", `YML file with format:
- path: /path1
  url:  url1
- path: /path2
  url:  url2
...
`)
	jsonFile := flag.String("json", "", `JSON file with format:
[
	{
		"path": "/path1",
		"url": "url1"
	},
	{
		"path": "/path2",
		"url": "url2"
	},
	...
]
	  

`)
	flag.Parse()
	return yamlFile, jsonFile
}

func readFile(path string) []byte {
	f, err := os.Open(path)
	handleError(err)

	fi, err := f.Stat()
	handleError(err)

	bytes := make([]byte, fi.Size())
	_, err = f.Read(bytes)
	handleError(err)

	return bytes
}

func main() {
	ymlFile, jsonFile := parseCli()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	handler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	var yamlBytes []byte
	if *ymlFile != "" {
		yamlBytes = readFile(*ymlFile)
	} else {
		yamlBytes = []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)
	}
	handler, err := urlshort.YAMLHandler(yamlBytes, handler)
	handleError(err)

	if *jsonFile != "" {
		jsonBytes := readFile(*jsonFile)
		handler, err = urlshort.JSONHandler(jsonBytes, handler)
		handleError(err)
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
