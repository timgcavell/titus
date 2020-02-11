package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var dir string

func main() {
	dir = "./"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	fs := http.FileServer(http.Dir(filepath.Join(dir + "static")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveTemplate)

	log.Printf("Serving %s\n", dir)
	http.ListenAndServe(":3000", nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	layoutPath := filepath.Join(dir, "templates", "layout.html")
	requestPath := r.URL.Path
	if requestPath == "/" {
		requestPath = "/index"
	}
	filePath := filepath.Join(dir, "templates", filepath.Clean(requestPath+".html"))

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	if fileInfo.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(layoutPath, filePath)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
