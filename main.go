package main

import (
	"net/http"
)

func main() {
	srvMux := http.NewServeMux()
	srvMux.Handle("/app/",
		http.StripPrefix(
			"/app/",
			http.FileServer(http.Dir("public_html"))))
	srvMux.HandleFunc("/healthz", healthHand)
	httpD := &http.Server{
		Handler: srvMux,
		Addr:    ":8080",
	}
	httpD.ListenAndServe()
}
