// Copyright 2013, Amahi.  All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

// server connection related functions

package spdy

import (
        "bytes"
        "fmt"
	"net/http"
)
//for hijacking
type Container struct {
	server      *http.Server
	chandler    *http.Handler
}
type buffer struct {
	bytes.Buffer
}

func (b *buffer) Close() error {
	return nil
}

// placeholder for proper error handling.
func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func(C *Container) hijack(w http.ResponseWriter, r *http.Request) {
        // re-purpose the connection.
	conn, _, err := w.(http.Hijacker).Hijack()
	handle(err)
	fmt.Println("re-purpose the connection")
        
        
	//respond to client
	buf := new(buffer)
        buf.WriteString("Hello from P")
	
	res := &http.Response{
		Status:        "200 Connection Established",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          buf,
		ContentLength: int64(buf.Len()),
	}
	handle(res.Write(conn))
	
	
        //change the handler
	C.server.Handler = *C.chandler
	if(C.server.Handler==nil) {
	        C.server.Handler = http.DefaultServeMux
	}
	fmt.Println("change the handler")
	
	//start session
	session := NewServerSession(conn, C.server)
	
	//serve
	fmt.Println("Serving started")
	handle(session.Serve())
	fmt.Println("Serving ended")
	
}

func ListenAndServe2(addr string, handler http.Handler) (err error) {
        //new container for hijacking
        C := new(Container)
        hServe := new(http.Server)
	mux := http.NewServeMux()
	mux.HandleFunc("/", C.hijack)
	fmt.Println("new container for hijacking")
	
	//save handler given by user
        C.chandler = &handler
        fmt.Println("save handler given by user")
        
        //save the server
        hServe.Handler = mux
	hServe.Addr = addr
	C.server = hServe
	fmt.Println("save the server")
	
	//start server hijack
        err = C.server.ListenAndServe()
        return err
}
