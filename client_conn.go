// Copyright 2013, Amahi.  All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

// client connection related functions

package spdy

import (
	"bytes"
	"net"
	"net/http"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
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

// Write always succeeds and writes to rw.Body.
func (rw *ResponseRecorder) Write(buf []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(200)
	}
	if rw.Body != nil {
		len, err := rw.Body.Write(buf)
		return len, err
	} else {
		rw.Body = new(bytes.Buffer)
		len, err := rw.Body.Write(buf)
		return len, err
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

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return &Client{}, err
	}
	return &Client{cn: conn}, nil
}

//to get a response from the client
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	session := NewClientSession(c.cn)
	go session.Serve()
	c.rr = new(ResponseRecorder)
	err := session.NewStreamProxy(req, c.rr)
	if err != nil {
		return &http.Response{}, err
	}
	resp := &http.Response{
		StatusCode:    c.rr.Code,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          &readCloser{c.rr.Body},
		ContentLength: int64(c.rr.Body.Len()),
	}
	return resp, nil
}
