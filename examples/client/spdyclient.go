package main

import (
        "github.com/nileshjagnik/spdy"
	"fmt"
	"io"
	"net/http"
)
func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
        client,err := spdy.NewClient("localhost:4040")
        handle(err)
        req,err := http.NewRequest("GET", "http://localhost:4040/banana", nil)
        handle(err)
        res,err := client.Do(req)
        handle(err)
        data := make([]byte, int(res.ContentLength))
        _,err = res.Body.(io.Reader).Read(data)
        fmt.Println(string(data))
        res.Body.Close()
}
