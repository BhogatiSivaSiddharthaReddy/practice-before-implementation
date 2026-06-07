package main

import (
	"fmt"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func hello(res http.ResponseWriter, req *http.Request) {

	panic("boom")

}

func TimerMiddlerware(han http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		han.ServeHTTP(w, r)

		fmt.Println(
			"Request took:",
			time.Since(start),
		)
	})
}

func AuthMiddleware(han http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		present := r.Header.Get("X-API-Key")

		if present == "" {
			http.Error(w, "Unauthorized", 401)
			return
		}

		han.ServeHTTP(w, r)
	})

}

func Chain(handler http.Handler, middleware ...Middleware) http.Handler {

	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}

	return handler

}

func RequestIDMiddleware(han http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fakeId := "12345"

		r.Header.Set("Request-Id", fakeId)

		han.ServeHTTP(w, r)
	})

}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request Started")

		getId := r.Header.Get("Request-Id")

		fmt.Println("Printing Request Id", getId)

		next.ServeHTTP(w, r)

		fmt.Println("request completed")
	})
}

func RecoveryMiddleware(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			err := recover()
			if err != nil {
				fmt.Println(
					"Recovered from panic:",
					err,
				)
				http.Error(w, "crashed", http.StatusInternalServerError)
			}

		}()

		next.ServeHTTP(w, r)

	})
}

func main() {
	fmt.Println("Day 8 practice")

	mid := Chain(http.HandlerFunc(hello), RecoveryMiddleware, RequestIDMiddleware, LoggingMiddleware, AuthMiddleware, TimerMiddlerware)
	http.Handle("/", mid)
	http.ListenAndServe(":8080", nil)

}
