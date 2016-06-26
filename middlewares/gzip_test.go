package middlewares

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {
	var mockHandler = func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"msg": "test message"}`)
		})
	}

	handler := mockHandler()
	handler = Gzip(handler)
	r, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Contains(t, w.Header().Get("Content-Encoding"), "gzip")

	reader, _ := gzip.NewReader(io.LimitReader(w.Body, 1<<20))
	b, _ := ioutil.ReadAll(reader)
	s := string(b[:])
	assert.Equal(t, `{"msg": "test message"}`, s, "unzip content compare")
}
