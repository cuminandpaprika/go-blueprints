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
	// get the room going
	go r.run()

	http.Handle("/room", r)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("/path/to/assets/"))))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: welcomeTemplate}))
	http.HandleFunc("/auth/", loginHandler)

	log.Printf("Serving webpage on %s", *hostNameAndPort)
	if err := http.ListenAndServe(*hostNameAndPort, nil); err != nil {
		log.Fatal(err)
	}

}
