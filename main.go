package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	//errFilesPath := GetEnvAsStringOrFallback("ERROR_FILES_PATH", "/www")

	t := template.Must(template.ParseGlob("./www/*.html"))

	http.HandleFunc("/", errorHandler(t))

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":8080"), nil); err != nil {
		panic(err)
	}
}
