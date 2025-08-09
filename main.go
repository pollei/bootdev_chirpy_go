package main

import (
	"net/http"
)

func main() {
	srvMux := http.NewServeMux()
	hcHand := newHandlerHitCounter(http.StripPrefix(
		"/app/",
		http.FileServer(http.Dir("public_html"))))
	srvMux.Handle("/app/", hcHand)
	srvMux.Handle("POST /reset", hcHand.reset)
	srvMux.Handle("GET /metrics", hcHand.metrics)
	srvMux.HandleFunc("GET /healthz", healthHand)
	httpD := &http.Server{
		Handler: srvMux,
		Addr:    ":8080",
	}
	httpD.ListenAndServe()
}
