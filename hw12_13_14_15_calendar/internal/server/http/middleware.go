package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}
func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		w, 
		http.StatusOK,
	}
}
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc { 
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
        start := time.Now()
		next.ServeHTTP(lrw, r)
		FileLog.Println(r.Host, fmt.Sprintf("[%v]", start.Format("02.01.2006 15:04:05")), r.Method, r.URL.Path, lrw.statusCode,
		time.Now().Sub(start), r.UserAgent())
	})
}
