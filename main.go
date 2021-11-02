package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
)

func main() {
	port := GetEnvAsStringOrFallback("PORT", "8080")
	errFilesPath := GetEnvAsStringOrFallback("ERROR_FILES_PATH", "/www")

	t := template.Must(template.ParseGlob(path.Join(errFilesPath, "*.html")))

	http.HandleFunc("/", errorHandler(t))

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
