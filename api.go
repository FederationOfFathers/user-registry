package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	nocache "github.com/rabeesh/negroni-nocache"
	"github.com/rs/cors"
	"github.com/uber-go/zap"
)

var router = mux.NewRouter()
var mwh = negroni.New()
var logger = zap.New(zap.NewJSONEncoder()).With(zap.String("module", "user-registry"))

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
	listenParts := strings.Split(listenOn, ":")
	listenPort, err := strconv.Atoi(listenParts[len(listenParts)-1])
	if err != nil {
		log.Fatalf("Unable to parse port from %s", listenOn)
	}
	err = consul.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name: "user-registry",
		Tags: []string{},
		Port: listenPort,
	})
	if err != nil {
		panic(err)
	}
	logger.Fatal(
		"error starting API http server",
		zap.String("listenOn", listenOn),
		zap.Error(http.ListenAndServe(listenOn, router)))
}
