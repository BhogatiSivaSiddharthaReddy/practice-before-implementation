package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// creating gateway servemux, handler and http server

	gateway := http.NewServeMux()

	gateway.HandleFunc("/", gatewayhandler)

	gatewayserver := &http.Server{
		Addr:    ":8080",
		Handler: gateway,
	}

	go func() {
		fmt.Println("gateway started on port 8080")
		if err := gatewayserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Gateway Error", err)
		}
	}()

	// creating backend servemux, handler, and http server
	backend := http.NewServeMux()

	backend.HandleFunc("/", backendhandler)

	backendserver := &http.Server{
		Addr:    ":9090",
		Handler: backend,
	}

	go func() {
		fmt.Println("backend started on port 9090")
		if err := backendserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("backend Error", err)
		}
	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch

	fmt.Println("shutting down gateway and backend")

}

func gatewayhandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("request reached gateway")

	NewRequest, err := http.NewRequest(req.Method, "http://localhost:9090", req.Body)

	if err != nil {
		fmt.Println("Failed to make Backend request")
		http.Error(res, "failed to make connection to backend", http.StatusInternalServerError)
		return
	}

	for head, values := range req.Header {
		for _, value := range values {
			NewRequest.Header.Add(head, value)
		}
	}

	NewRequest.Header.Add("X-Forwarded-For", req.RemoteAddr)
	client := &http.Client{}

	newres, err := client.Do(NewRequest)

	if err != nil {
		http.Error(res, "Fialed to make connection", http.StatusBadGateway)
		return
	}
	newres.Body.Close()

	for head, values := range newres.Header {
		for _, value := range values {
			res.Header().Add(head, value)
		}
	}

	// Set status code res.WriteHeader(resp.StatusCode) // Copy body io.Copy(res, resp.Body)
	res.WriteHeader(newres.StatusCode)

	io.Copy(res, newres.Body)
}

func backendhandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Backend recieved request")
	clientIP := req.Header.Get("X-Forwarded-For")
	fmt.Println("X-Forwarded-For:", clientIP)
	res.Header().Add(
		"X-Backend-Service",
		"users-api",
	)

	res.WriteHeader(http.StatusCreated)

	fmt.Fprintln(res, "Backend response received")
}
