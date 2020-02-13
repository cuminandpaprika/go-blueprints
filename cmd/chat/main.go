package main

import (
	"log"
	"net/http"
)

func handleWeb(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(` 
	  <html> 
		<head> 
		  <title>Chat</title> 
		</head> 
		<body> 
		  Let's chat! 
		</body> 
	  </html>`))
}

func main() {
	const hostNameAndPort string = ":8080"
	http.HandleFunc("/", handleWeb)
	log.Printf("Serving webpage on %s", hostNameAndPort)
	if err := http.ListenAndServe(hostNameAndPort, nil); err != nil {
		log.Fatal(err)
	}

}
