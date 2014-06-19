// Copyright 2013, Amahi.  All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

// lower level test functions

package spdy

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"testing"
	"time"
)

const HOST_PORT_API = "localhost:1443"
const HOST_PORT_SERVERS = "localhost:1444"
const HOST_PORT = "localhost:1444"
const SERVER_ROOT = "../testdata"

type handler struct {
	data []byte
	rt   string
}

type stats_s struct {
	sync.Mutex
	incoming int
	serving  int
}

var stats stats_s

// Used in sending the response.
// Essentially, this is just adding
// the Close method so that it fulfils
// the io.ReadCloser interface.
type buffer struct {
	bytes.Buffer
}

func (b *buffer) Close() error {
	return nil
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	if rq.Body != nil {
		h.data = make([]byte, int(rq.ContentLength))
		_, err := rq.Body.(io.Reader).Read(h.data)
		if err != nil {
			fmt.Println(err)
		}
		filename := "/tmp/postdat"
		f, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)
		}
		n, err := f.Write(h.data)
		if err != nil {
			fmt.Println(n, err)
		}
		f.Close()
	}
	fileserver := http.FileServer(http.Dir(h.rt))
	fileserver.ServeHTTP(rw, rq)
}

type Proxy struct {
	session *Session
}

func (p *Proxy) RequestFromC(w http.ResponseWriter, r *http.Request) error {
	if p.session == nil {
		log.Println("Warning: Could not serve request because C is not connected.")
		http.NotFound(w, r)
		return nil
	}

	u := r.URL
	if u.Host == "" {
		u.Host = HOST_PORT_API
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	err := p.session.NewStreamProxy(r, w)
	return err
}

func (p *Proxy) ServeC(w http.ResponseWriter, r *http.Request) {
	// clean up the old connection
	if p.session != nil {
		p.session.Close()
	}

	// Read in the request body.
	buf := new(buffer)
	_, err := io.Copy(buf, r.Body)
	handle(err)
	handle(r.Body.Close())

	// re-purpose the connection.
	conn, _, err := w.(http.Hijacker).Hijack()
	handle(err)

	// send the 200 to C.
	buf.Reset()
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

	// prepare for serving requests from A.
	session := NewClientSession(conn)
	p.session = session
	session.Serve()
}

func (p *Proxy) ServeA(w http.ResponseWriter, r *http.Request) {
	err := p.RequestFromC(w, r)
	if err != nil {
		log.Println("Encountered an error serving API request:", err)
	}
}

func (p *Proxy) DebugURL(w http.ResponseWriter, r *http.Request) {
	stats.Lock()
	fmt.Fprintf(w, "goroutines:  %d\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "incoming: %d\nserving: %d\n", stats.incoming, stats.serving)
	stats.Unlock()
}

func p() {
	certFile := "cert/cert.pem"
	keyFile := "cert/cert.key"
	proxy := new(Proxy)
	http.HandleFunc("/", proxy.ServeC)

	go func() { // Serve C
		err := http.ListenAndServeTLS(HOST_PORT_SERVERS, certFile, keyFile, nil)
		handle(err)
	}()

	hServe := new(http.Server)
	mux := http.NewServeMux()
	mux.HandleFunc("/", proxy.ServeA)
	mux.HandleFunc("/debug", proxy.DebugURL)
	hServe.Handler = mux
	hServe.Addr = HOST_PORT_API
	handle(hServe.ListenAndServe()) // Serve H
}

func c() {

	root := SERVER_ROOT
	for {
		const SLEEP_RETRY = 5
		var conn *tls.Conn
		var err error
		for i := 0; i < 10; i++ {
			// connect to P.
			conn, err = tls.Dial("tcp", HOST_PORT, &tls.Config{InsecureSkipVerify: true})
			if err != nil {
				time.Sleep(100 * time.Millisecond)
			} else {
				break
			}
		}
		if conn == nil {
			log.Println("Failed to connect. Waiting", SLEEP_RETRY, "seconds.")
			time.Sleep(SLEEP_RETRY * time.Second)
			continue
		}

		// build the request
		buf := new(bytes.Buffer)
		_, err = buf.WriteString("Hello from C")
		handle(err)
		req, err := http.NewRequest("PUT", "https://"+HOST_PORT, buf)
		handle(err)

		// make the client connection
		client := httputil.NewClientConn(conn, nil)
		res, err := client.Do(req)
		if err != nil {
			log.Println("Error: Failed to make connection to P:", err)
			continue
		}
		buf.Reset()
		_, err = io.Copy(buf, res.Body)
		handle(err)

		c, _ := client.Hijack()
		conn = c.(*tls.Conn)
		server := new(http.Server)
		server.Handler = &handler{data: nil, rt: root}
		session := NewServerSession(conn, server)
		session.Serve()
	}
}
func init() {
	go p()
	go c()
	time.Sleep(time.Second)
	handle(os.Chdir("./integration-tests"))
	SetLog(ioutil.Discard)
}
func ServerTestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func TestFrames(t *testing.T) {
	//make server
	mux := http.NewServeMux()
	mux.HandleFunc("/", ServerTestHandler)
	server := &Server{
		Addr:    "localhost:4040",
		Handler: mux,
	}
	go server.ListenAndServe()
	time.Sleep(200 * time.Millisecond)

	//make client
	client, err := NewClient("localhost:4040")
	if err != nil {
		t.Fatal(err.Error())
	}

	//now send requests and test
	req, err := http.NewRequest("GET", "http://localhost:4040/banana", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	data := make([]byte, int(res.ContentLength))
	_, err = res.Body.(io.Reader).Read(data)
	if string(data) != "Hi there, I love banana!" {
		t.Fatal("Unexpected Data")
	}
	res.Body.Close()

	//another request
	req, err = http.NewRequest("POST", "http://localhost:4040/monkeys", bytes.NewBufferString("hello=world"))
	if err != nil {
		t.Fatal(err.Error())
	}

	res, err = client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	data = make([]byte, int(res.ContentLength))
	_, err = res.Body.(io.Reader).Read(data)
	if string(data) != "Hi there, I love monkeys!" {
		t.Fatal("Unexpected Data")
	}
	res.Body.Close()

	//settings frame test
	set := new(settings)
	var svp []settingsValuePairs
	svp = append(svp, settingsValuePairs{flags: 0, id: 4, value: 6})   //set SETTINGS_MAX_CONCURRENT_STREAMS to 6
	svp = append(svp, settingsValuePairs{flags: 0, id: 3, value: 400}) //set SETTINGS_ROUND_TRIP_TIME to 400ms
	set.flags = 0
	set.count = 2
	set.svp = svp
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, &set.count)
	if err != nil {
		t.Fatal(err.Error())
	}

	for i := uint32(0); i < set.count; i++ {
		err = binary.Write(buf, binary.BigEndian, &set.svp[i].flags)
		if err != nil {
			t.Fatal(err.Error())
		}
		err = binary.Write(buf, binary.BigEndian, &set.svp[i].id)
		if err != nil {
			t.Fatal(err.Error())
		}
		err = binary.Write(buf, binary.BigEndian, &set.svp[i].value)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	settings_frame := controlFrame{kind: FRAME_SETTINGS, flags: 0, data: buf.Bytes()}
	client.ss.out <- settings_frame
	time.Sleep(200 * time.Millisecond)

	//rstStreamtest - first, start stream on client
	str := client.ss.NewClientStream()
	if str == nil {
		t.Fatal("ERROR in NewClientStream: cannot create stream")
		return
	}
	str.sendRstStream()

	//ping test
	ping, err := client.Ping(time.Second)
	if err != nil {
		t.Fatal(err.Error())
	}

	if ping == false {
		t.Fatal("Unable to ping server from client")
	}

	//close client
	err = client.Close()
	if err != nil {
		t.Fatal(err.Error())
	}
	//server close
	server.Close()
}

func TestGet(t *testing.T) {
	cmd := exec.Command("bash", "test-01-basic-root-dir-listing.sh")
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err.Error())
	}
	outstr := string(out)
	result := "Head: PASS\nBody: FAIL\n"
	if result != outstr {
		fmt.Println(outstr)
		t.Fatal("Unexpected Output")
	}
}

func TestImage(t *testing.T) {
	cmd := exec.Command("bash", "test-02-image.sh")
	out, err := cmd.Output()
	if err != nil {
		t.Error(err.Error())
	}
	outstr := string(out)
	result := "Head: PASS\nBody: PASS\n"
	if result != outstr {
		t.Error("Unexpected Output")
	}
}

func TestVideoAVI(t *testing.T) {
	cmd := exec.Command("bash", "test-03-video-avi.sh")
	out, err := cmd.Output()
	if err != nil {
		t.Error(err.Error())
	}
	outstr := string(out)
	result := "Head: PASS\nBody: PASS\n"
	if result != outstr {
		t.Error("Unexpected Output")
	}
}

func TestVideoMKV(t *testing.T) {
	cmd := exec.Command("bash", "test-04-video-mkv.sh")
	out, err := cmd.Output()
	if err != nil {
		t.Error(err.Error())
	}
	outstr := string(out)
	result := "Head: FAIL\nBody: PASS\n"
	if result != outstr {
		t.Error("Unexpected Output")
	}
}

func TestRootWithIF(t *testing.T) {
	cmd := exec.Command("bash", "test-06-root-with-if-modified.sh")
	out, err := cmd.Output()
	if err != nil {
		t.Error(err.Error())
	}
	outstr := string(out)
	result := "Head: FAIL\n"
	if result != outstr {
		t.Error("Unexpected Output")
	}
}

func TestRangeReq(t *testing.T) {
	cmd := exec.Command("bash", "test-07-range-request.sh")
	out, err := cmd.Output()
	if err != nil {
		t.Error(err.Error())
	}
	outstr := string(out)
	result := "Head: PASS\nBody: PASS\n"
	if result != outstr {
		t.Error("Unexpected Output")
	}
}

func TestMoviePlaySafari(t *testing.T) {
	cmd := exec.Command("bash", "test-08-movie-play-in-safari.sh")
	out, err := cmd.Output()
	if err != nil {
		t.Error(err.Error())
	}
	outstr := string(out)
	result := "Head: PASS\nBody: PASS\n"
	if result != outstr {
		t.Error("Unexpected Output")
	}
}

func TestBasicPost(t *testing.T) {
	cmd := exec.Command("bash", "test-100-basic-post.sh")
	out, err := cmd.Output()
	if err != nil {
		t.Error(err.Error())
	}
	outstr := string(out)
	result := "Body: PASS\nData Receive: PASS\n"
	if result != outstr {
		t.Error("Unexpected Output")
	}
}

func TestHeadReq(t *testing.T) {
	cmd := exec.Command("bash", "test-101-head-req.sh")
	out, err := cmd.Output()
	if err != nil {
		t.Error(err.Error())
	}
	outstr := string(out)
	result := "PASS\n"
	if result != outstr {
		t.Error("Unexpected Output")
	}
}

/*
func TestSlowCall(t *testing.T) {
        cmd := exec.Command("bash","test-80-sloow-call.sh")
        out,err := cmd.Output()
        if err != nil {
                t.Error(err.Error())
        }
        outstr := string(out)
        result := "Head: PASS\nBody: PASS\n"
        if result != outstr {
                t.Error("Unexpected Output")
        }
}
*/
