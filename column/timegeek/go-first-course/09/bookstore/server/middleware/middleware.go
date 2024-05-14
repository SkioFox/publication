package middleware

import (
	"fmt"
	"log"
	"mime"
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//  recv a POST request from 127.0.0.1:63884
		log.Printf("recv a %s request from %s", req.Method, req.RemoteAddr)
		// req: &{POST /book HTTP/1.1 1 1 map[Accept:[*/*] Accept-Encoding:[gzip, deflate, br] Cache-Control:[no-cache] Connection:[keep-alive] Content-Length:[89] Content-Type:[application/json] Postman-Token:[19a99b8a-4334-4a2f-86a4-99a5b9cf4e89] User-Agent:[PostmanRuntime/7.37.0]] 0xc0000de6c0 <nil> 89 [] false 127.0.0.1:8080 map[] map[] <nil> map[] 127.0.0.1:63884 /book <nil> <nil> <nil> 0xc0000fa4b0} &{0xc00009c288 <nil> <nil> false true {0 0} false false false 0x139e920} /book
		fmt.Println("req:", req, req.Body, req.URL)
		next.ServeHTTP(w, req)
	})
}

func Validating(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		contentType := req.Header.Get("Content-Type")
		mediatype, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if mediatype != "application/json" {
			http.Error(w, "invalid Content-Type", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, req)
	})
}
