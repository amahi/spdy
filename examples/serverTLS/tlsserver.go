package main

import (
	"fmt"
	"github.com/amahi/spdy"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
        http.HandleFunc("/", handler)
	spdy.ListenAndServeTLS("localhost:4040", "server.pem", "server.key" , nil)
}
