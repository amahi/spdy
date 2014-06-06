package main

import (
        "bytes"
	"fmt"
	"net"
	"net/http"
	//"net/http/httputil"
	"github.com/amahi/spdy"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}


// ResponseRecorder is an implementation of http.ResponseWriter that
// records its mutations for later inspection in tests.
type ResponseRecorder struct {
	Code      int           // the HTTP response code from WriteHeader
	HeaderMap http.Header   // the HTTP response headers
	Body      *bytes.Buffer // if non-nil, the bytes.Buffer to append written data to
	Flushed   bool

	wroteHeader bool
}

// NewRecorder returns an initialized ResponseRecorder.
func NewRecorder() *ResponseRecorder {
	return &ResponseRecorder{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Code:      200,
	}
}

// Header returns the response headers.
func (rw *ResponseRecorder) Header() http.Header {
	m := rw.HeaderMap
	if m == nil {
		m = make(http.Header)
		rw.HeaderMap = m
	}
	return m
}

// Write always succeeds and writes to rw.Body, if not nil.
func (rw *ResponseRecorder) Write(buf []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(200)
	}
	if rw.Body != nil {
		rw.Body.Write(buf)
	}
	return len(buf), nil
}

// WriteHeader sets rw.Code.
func (rw *ResponseRecorder) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.Code = code
	}
	rw.wroteHeader = true
}

// Flush sets rw.Flushed to true.
func (rw *ResponseRecorder) Flush() {
	if !rw.wroteHeader {
		rw.WriteHeader(200)
	}
	rw.Flushed = true
}

func main() {
        req, err := http.NewRequest("GET", "http://localhost:4040/hello",nil)
	handle(err)
        conn, err := net.Dial("tcp", "localhost:4040")
        handle(err)
        spdy.EnableDebug()
        // make the client connection
        
	//client := httputil.NewClientConn(conn, nil)
	
	/*
	//send hello
        res, err := client.Do(req)
        fmt.Println(res.Status)
        */
        //hijack for spdy
	//conn, _ = client.Hijack()
	session := spdy.NewClientSession(conn)
	fmt.Println("Ready")
	go session.Serve()
	//res, err = client.Do(req)
	w := new(ResponseRecorder)
	fmt.Println("Serving")
	handle(session.NewStreamProxy(req, w))
	fmt.Println(w.Body)
}
