package middleware

import (
	"log"
	"mime"
	"net/http"
)

// Logging 中间件
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// 打印请求信息
		log.Printf("recv a %s request from %s", req.Method, req.RemoteAddr)
		next.ServeHTTP(w, req)
	})
}

// Validating 中间件
// 验证请求的 Content-Type
func Validating(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		contentType := req.Header.Get("Content-Type")
		mediatype, _, err := mime.ParseMediaType(contentType)
		// 如果解析出错，则返回错误
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// 如果Content-Type不是application/json，则返回错误
		if mediatype != "application/json" {
			http.Error(w, "invalid Content-Type", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, req)
	})
}
