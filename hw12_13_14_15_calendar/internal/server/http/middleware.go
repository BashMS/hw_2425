package internalhttp

import (
	"fmt"
	"net/http"
	"time"
	"net"
	"strings"
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

func getIP(r *http.Request) string {
    //Get IP from the X-REAL-IP header
    ip := r.Header.Get("X-REAL-IP")
    netIP := net.ParseIP(ip)
    if netIP != nil {
        return ip
    }

    //Get IP from X-FORWARDED-FOR header
    ips := r.Header.Get("X-FORWARDED-FOR")
    splitIps := strings.Split(ips, ",")
    for _, ip := range splitIps {
        netIP := net.ParseIP(ip)
        if netIP != nil {
            return ip
        }
    }

    //Get IP from RemoteAddr
    ip, _, _ = net.SplitHostPort(r.RemoteAddr)
    netIP = net.ParseIP(ip)
    if netIP != nil {
        return ip
    }
    return ""
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc { 
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		start := time.Now()
		next.ServeHTTP(lrw, r)
		FileLog.Println(getIP(r), fmt.Sprintf("[%v]", start.Format("02.01.2006 15:04:05")), r.Method, r.URL.Path, lrw.statusCode,
		time.Now().Sub(start), r.UserAgent())
	})
}
