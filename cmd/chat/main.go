package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/cuminandpaprika/go-blueprints/pkg/trace"
)

// templateHandler represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

const (
	templateDir     = "templates"
	welcomeTemplate = "chat.html"
	hostNameAndPort = ":8080"
)

// serveHTTP handles HTTP requests with lazy loading
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(
		func() {
			t.templ = template.Must(template.ParseFiles(filepath.Join(templateDir, t.filename)))

		})
	t.templ.Execute(w, r)
}

func main() {
	var hostNameAndPort = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/room", r)
	// get the room going
	go r.run()

	handler := templateHandler{filename: welcomeTemplate}
	http.HandleFunc("/", handler.ServeHTTP)

	log.Printf("Serving webpage on %s", *hostNameAndPort)
	if err := http.ListenAndServe(*hostNameAndPort, nil); err != nil {
		log.Fatal(err)
	}

}
