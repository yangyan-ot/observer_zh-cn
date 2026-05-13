package socket

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type noopResponseWriter struct {
	header http.Header
	code   int
}

func (n *noopResponseWriter) Header() http.Header         { return n.header }
func (n *noopResponseWriter) Write(b []byte) (int, error) { return len(b), nil }
func (n *noopResponseWriter) WriteHeader(code int)        { n.code = code }

func newTokenValidator(middlewareFn gin.HandlerFunc) func(string) bool {
	return func(tokenStr string) bool {
		w := &noopResponseWriter{header: make(http.Header)}
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenStr)
		c.Request = req
		middlewareFn(c)
		return !c.IsAborted()
	}
}
