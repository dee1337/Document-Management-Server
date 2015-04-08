package api

import (
	"net/http"
	//"net/http/httptest"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

/*
func PathWorks(t *testing.T) {
	http.HandleFunc("/list", Files)
	http.ListenAndServe(":8080", nil)
	resp, err := http.Get("http://localhost:8080/list")
	defer resp.Close()
    httptest.Server
	if err != nil {
		t.Error("Got an error")
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	// Test that status code was 200 etc.
}


func TestHome(t *testing.T) {
	//mockDb := MockDb{}
	homeHandle := Files
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	homeHandle.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}
*/
func FileHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	io.WriteString(w, "START\n")
	f := w.(http.Flusher)
	f.Flush()
	time.Sleep(200 * time.Millisecond)
	io.WriteString(w, "WORKING\n")
	f.Flush()
	time.Sleep(200 * time.Millisecond)
	io.WriteString(w, "DONE\n")
	return
}

func setupMockServer(t *testing.T) {
	http.HandleFunc("/list", FileHandler)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatalf("failed to listen - %s", err.Error())
	}
	go func() {
		err = http.Serve(ln, nil)
		if err != nil {
			t.Fatalf("failed to start HTTP server - %s", err.Error())
		}
	}()
	addr = ln.Addr()
}
