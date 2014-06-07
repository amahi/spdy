// Copyright 2013, Amahi.  All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

// server connection related functions

package spdy

import (
	"net"
	"net/http"
	"time"
)

func handleConnection(conn net.Conn, addr string, handler http.Handler) {
	hserve := new(http.Server)
	if handler == nil {
		hserve.Handler = http.DefaultServeMux
	} else {
		hserve.Handler = handler
	}
	hserve.Addr = addr
	session := NewServerSession(conn, hserve)
	handle(session.Serve())
}

func ListenAndServe(addr string, handler http.Handler) (err error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// handle error
	}
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("http: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		go handleConnection(conn, addr, handler)
	}
}
