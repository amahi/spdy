// Copyright 2013, Amahi.  All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

// Debug and logging related functions

package spdy

import (
	"io"
	"io/ioutil"
	logging "log"
	"os"
)

// regular app logging - enabled by default
var log = logging.New(os.Stderr, "[SPDY] ", logging.LstdFlags|logging.Lshortfile)

// app logging for the purposes of debugging - disabled by default
var debug = logging.New(ioutil.Discard, "[SPDY DEBUG] ", logging.LstdFlags)

func EnableDebug() {
	debug = logging.New(os.Stdout, "[SPDY DEBUG] ", logging.LstdFlags)
}

func SetLog(w io.Writer) {
	log = logging.New(w, "[SPDY] ", logging.LstdFlags|logging.Lshortfile)
}
