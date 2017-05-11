package testsupport

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

type TestHandler struct {
	nBytesReceived    int
	nRequestsInFlight int64
}

func (th *TestHandler) NumRequestsInFlight() int {
	return int(atomic.LoadInt64(&th.nRequestsInFlight))
}

func (th *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&th.nRequestsInFlight, 1)
	fmt.Println("handler start")

	defer func() {
		atomic.AddInt64(&th.nRequestsInFlight, -1)
		fmt.Println("handler end")
	}()

	inputBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(bytes.ToUpper(inputBytes)); err != nil {
		panic(err)
	}
}

type SingleByteWriter struct {
	Destination   io.Writer
	ByteInspector func(b byte) error
}

func (w *SingleByteWriter) Write(buffer []byte) (int, error) {
	total := 0
	for i, b := range buffer {
		if w.ByteInspector != nil {
			if err := w.ByteInspector(b); err != nil {
				return total, err
			}
		}

		n, err := w.Destination.Write(buffer[i : i+1])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}
