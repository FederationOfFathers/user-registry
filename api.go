package main

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	nocache "github.com/rabeesh/negroni-nocache"
	"github.com/rs/cors"
	"github.com/uber-go/zap"
)

var router = mux.NewRouter()
var mwh = negroni.New()
var logger = zap.NewJSON().With(zap.String("module", "user-registry"))

var cmw = cors.New(cors.Options{
	AllowedOrigins: []string{
		"http://127.0.0.1",
		"http://localhost",
	},
	AllowedMethods: []string{
		"GET",
		"PUT",
	},
	AllowCredentials: true,
})

func mw(fn func(w http.ResponseWriter, r *http.Request)) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.*"},
		AllowCredentials: true,
	})
	return gziphandler.GzipHandler(
		c.Handler(
			negroni.New(
				&httpLogger{},
				negroni.NewRecovery(),
				nocache.New(true),
				negroni.Wrap(
					http.HandlerFunc(fn),
				),
			),
		),
	)
}

func mindHTTP() {
	logger.Fatal(
		"error starting API http server",
		zap.String("listenOn", listenOn),
		zap.Error(http.ListenAndServe(listenOn, router)))
}
