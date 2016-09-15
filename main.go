package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type key int

const requestIDKey key = 0

func main() {
	addr := flag.String("addr", ":8080", "ip:port to listen on")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/hello", helloHandler).Methods("GET", "HEAD")

	if err := http.ListenAndServe(*addr, middleware(router)); err != nil {
		panic(err)
	}
}

func helloHandler(rw http.ResponseWriter, req *http.Request) {
	reqID := requestIDFromContext(req.Context())
	fmt.Fprintf(rw, "Hello request ID %v\n", reqID)
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := newContextWithRequestID(req.Context(), req)
		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}

func newContextWithRequestID(
	ctx context.Context,
	req *http.Request) context.Context {

	reqID := req.Header.Get("X-Request-ID")
	if reqID == "" {
		reqID = generateRandomID()
	}

	return context.WithValue(ctx, requestIDKey, reqID)
}

func requestIDFromContext(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)
}

func generateRandomID() string {
	return uuid.NewV4().String()
}
