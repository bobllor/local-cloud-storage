package api

import (
	"fmt"
	"net/http"
	"testing"
)

type MockWriter struct {
}

func (mw *MockWriter) Header() http.Header {
	return make(http.Header)
}

func (mw *MockWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (mw *MockWriter) WriteHeader(statusCode int) {
	fmt.Println(statusCode)
}

func TestPostRegisterUser(t *testing.T) {
}
