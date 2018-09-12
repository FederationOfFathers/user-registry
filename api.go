package main

import (
	"net/http"

	"github.com/FederationOfFathers/consul"
	"github.com/NYTimes/gziphandler"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	nocache "github.com/rabeesh/negroni-nocache"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

var router = mux.NewRouter()
var mwh = negroni.New()
var logger *zap.Logger

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

func init() {
	l, _ := zap.NewProduction()
	logger = l.With(zap.String("module", "user-registry"))
}

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
	if err := consul.RegisterOn("user-registry", listenOn); err != nil {
		panic(err)
	}
	logger.Fatal(
		"error starting API http server",
		zap.String("listenOn", listenOn),
		zap.Error(http.ListenAndServe(listenOn, router)))
}
