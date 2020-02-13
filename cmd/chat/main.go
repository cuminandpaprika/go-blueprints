package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// templateHandler represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

const templateDir string = "templates"
const welcomeTemplate string = "chat.html"
const hostNameAndPort string = ":8080"

// serveHTTP handles HTTP requests with lazy loading
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(
		func() {
			t.templ = template.Must(template.ParseFiles(filepath.Join(templateDir, t.filename)))

		})
	t.templ.Execute(w, nil)
}

func main() {
	handler := templateHandler{filename: welcomeTemplate}
	http.HandleFunc("/", handler.ServeHTTP)

	log.Printf("Serving webpage on %s", hostNameAndPort)
	if err := http.ListenAndServe(hostNameAndPort, nil); err != nil {
		log.Fatal(err)
	}

}
