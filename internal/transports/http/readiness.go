package http

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

var sigtermCalled = false

func (h HttpServer) initSigtermListener(errs chan error) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		err := <-c
		sigtermCalled = true
		errs <- fmt.Errorf("%s", err)
	}()
}

func (h HttpServer) registerReadinessRoute(rtr *mux.Router) {
	handler := http.HandlerFunc(makeReadinessHTTPHandler())
	rtr.Methods("GET").
		Path("/readiness").
		Handler(handler)
}

func makeReadinessHTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if sigtermCalled {
			w.WriteHeader(503)
			w.Write([]byte("{\"ready\": false}"))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("{\"ready\": true}"))
	}
}
