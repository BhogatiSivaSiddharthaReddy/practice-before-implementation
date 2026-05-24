package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func gatewayHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Request reached gateway")

	// Create new request to backend
	backendReq, err := http.NewRequest(
		req.Method,
		"http://localhost:8081",
		req.Body,
	)
	if err != nil {
		http.Error(res, "Failed to create backend request", http.StatusInternalServerError)
		return
	}

	// Copy headers from original request
	for key, values := range req.Header {
		for _, value := range values {
			backendReq.Header.Add(key, value)
		}
	}

	// Add X-Forwarded-For header
	backendReq.Header.Add("X-Forwarded-For", req.RemoteAddr)

	client := &http.Client{}

	resp, err := client.Do(backendReq)
	if err != nil {
		http.Error(res, "Failed to reach backend", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy backend headers to response
	for key, values := range resp.Header {
		for _, value := range values {
			res.Header().Add(key, value)
		}
	}

	// Set status code
	res.WriteHeader(resp.StatusCode)

	// Copy body
	io.Copy(res, resp.Body)
}

func backendHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Request reached backend")
	clientIP := req.Header.Get("X-Forwarded-For")
	fmt.Println("X-Forwarded-For:", clientIP)
	fmt.Fprintln(res, "Backend response received")
}

func main() {
	// Backend server
	backendMux := http.NewServeMux()
	backendMux.HandleFunc("/", backendHandler)

	backend := &http.Server{
		Addr:    ":8081",
		Handler: backendMux,
	}

	go func() {
		fmt.Println("Backend running on :8081")
		if err := backend.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Backend error:", err)
		}
	}()

	// Gateway server
	gatewayMux := http.NewServeMux()
	gatewayMux.HandleFunc("/", gatewayHandler)

	gateway := &http.Server{
		Addr:    ":8080",
		Handler: gatewayMux,
	}

	go func() {
		fmt.Println("Gateway running on :8080")
		if err := gateway.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Gateway error:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down servers...")
}
