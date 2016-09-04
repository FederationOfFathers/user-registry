package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/rs/cors"
)

var router = mux.NewRouter()
var mwh = negroni.New()

func mindHTTP() {
	// CORS
	mwh.Use(cors.New(cors.Options{AllowedOrigins: []string{"fofgatming", "127.0.0.1", "localhost"}}))
	// Recovery
	mwh.Use(negroni.NewRecovery())
	// Logging TODO: Replace
	mwh.Use(negroni.NewLogger())
	// Gzip Compression of Responses
	mwh.Use(gzip.Gzip(gzip.DefaultCompression))
	// Wrap the router
	mwh.UseHandler(router)
	// go TODO: replace
	mwh.Run(listenOn)
}
