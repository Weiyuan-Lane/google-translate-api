package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/weiyuan-lane/google-translate-api/internal/services/googletranslatewrapper"
	"github.com/weiyuan-lane/google-translate-api/internal/transports/http/services/googletranslate"
	loggerutils "github.com/weiyuan-lane/google-translate-api/internal/utils/logger"
)

type HttpServer struct {
	LivelinessProbePort      string
	Port                     string
	Logger                   *loggerutils.Logger
	GracefulShutdownSeconds  int
	EnableHTTP2              bool
	GoogleTranslateV2Wrapper googletranslatewrapper.TranslateV2Wrapper
	GoogleTranslateV3Wrapper googletranslatewrapper.TranslateV3Wrapper
}

func (h HttpServer) ListenAndServe() {
	h.initLivelinessHTTPServer()
	h.initCoreHTTPServer()
}

func (h HttpServer) initLivelinessHTTPServer() {
	address := ":" + h.LivelinessProbePort

	livenessHTTPHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("{}"))
	})

	go func() {
		server := &http.Server{
			Addr:    address,
			Handler: livenessHTTPHandler,
		}
		server.ListenAndServe()
	}()
}

func (h HttpServer) initCoreHTTPServer() {
	router := mux.NewRouter()
	errs := make(chan error)
	h.initSigtermListener(errs)
	h.registerReadinessRoute(router)
	h.registerServices(router)
	handler := h.makeCORSWrappedHTTPHandler(router)
	address := ":" + h.Port
	server := h.makeHttpServerFrom(address, handler)

	go func() {
		errs <- server.ListenAndServe()
	}()

	h.Logger.Info("Serving from port " + h.Port)
	h.Logger.Info((<-errs).Error())

	gracefulShutdownTime := time.Duration(h.GracefulShutdownSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTime)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		h.Logger.Info(fmt.Sprintf("Server Shutdown Failed:%+v", err))
	} else {
		h.Logger.Info("Graceful shutdown completed")
	}
}

func (h HttpServer) registerServices(router *mux.Router) {
	googleTranslateSvc := googletranslate.GoogleTranslateService{
		Logger:             h.Logger,
		TranslateV2Wrapper: h.GoogleTranslateV2Wrapper,
		TranslateV3Wrapper: h.GoogleTranslateV3Wrapper,
	}

	h.registerRoutes(
		router,
		googleTranslateSvc,
	)
}

func (h HttpServer) makeHttpServerFrom(address string, handler http.Handler) *http.Server {
	var server *http.Server

	if h.EnableHTTP2 {
		h2s := &http2.Server{}
		server = &http.Server{
			Addr:    address,
			Handler: h2c.NewHandler(handler, h2s),
		}
	} else {
		server = &http.Server{
			Addr:    address,
			Handler: handler,
		}
	}

	return server
}
