package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//
// BACKEND
//

func backendHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Backend received path:", r.URL.Path)

	fmt.Fprintln(w, "response from backend")
}

//
// MAIN
//

func main() {

	//
	// Backend Server
	//

	backendMux := http.NewServeMux()

	backendMux.HandleFunc("/", backendHandler)

	go func() {

		fmt.Println("Backend running on :8081")

		http.ListenAndServe(":8081", backendMux)
	}()

	//
	// Gateway Proxy
	//

	target, _ := url.Parse("http://localhost:8081")

	proxy := httputil.NewSingleHostReverseProxy(target)

	fmt.Println("Gateway running on :8080")

	http.ListenAndServe(":8080", proxy)
}
