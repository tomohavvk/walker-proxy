package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func verifyRequest(req *http.Request) bool {

	return true
}

func proxyHandler(target *url.URL) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req)
		if verifyRequest(req) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host

			proxy := httputil.NewSingleHostReverseProxy(target)
			proxy.ServeHTTP(w, req)
		} else {
			http.Error(w, "Some forbidden message...", http.StatusForbidden)
		}
	})
}

func main() {
	targetURL, err := url.Parse("http://localhost:80")
	if err != nil {
		fmt.Println("Error parsing target URL:", err)
		return
	}

	proxy := proxyHandler(targetURL)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		proxy.ServeHTTP(w, req)
	})

	fmt.Println("Proxy server listening on :8888...")
	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
