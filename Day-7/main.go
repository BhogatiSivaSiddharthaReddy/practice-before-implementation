package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

//
// BACKEND
//

func backendHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Path:", r.URL.Path)

	fmt.Println(
		"Gateway Version:",
		r.Header.Get("X-Gateway-Version"),
	)

	fmt.Println(
		"X-Forwarded-For:",
		r.Header.Get("X-Forwarded-For"),
	)

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

	proxy.Director = nil

	proxy.Rewrite = func(pr *httputil.ProxyRequest) {
		pr.SetURL(target)
		pr.SetXForwarded()

		pr.Out.Header.Set(
			"X-Gateway-Version",
			"v1",
		)

		pr.Out.URL.Path = strings.TrimPrefix(
			pr.In.URL.Path,
			"/api",
		)

	}

	proxy.ModifyResponse = func(r *http.Response) error {
		fmt.Println(
			"Backend returned:",
			r.StatusCode,
		)

		r.Header.Set(
			"X-Processed-By",
			"MyGateway",
		)
		return nil
	}

	proxy.ErrorHandler = func(
		w http.ResponseWriter,
		r *http.Request,
		err error,
	) {

		fmt.Printf(
			"%s %s failed: %v\n",
			r.Method,
			r.URL.Path,
			err,
		)

		http.Error(
			w,
			"Backend unavailable",
			503,
		)
	}

	fmt.Println("Gateway running on :8080")

	http.ListenAndServe(":8080", proxy)
}
