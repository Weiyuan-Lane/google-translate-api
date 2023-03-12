package http

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/weiyuan-lane/google-translate-api/internal/transports/http/services/googletranslate"
)

func (h HttpServer) registerRoutes(
	rtr *mux.Router,
	googleTranslateService googletranslate.GoogleTranslateService,
) {

	rtr.Methods("POST").Path("/google-translate/v2/translate").Handler(googleTranslateService.GoogleTranslateV2TranslateHandler())
	rtr.Methods("POST").Path("/google-translate/v2/detect").Handler(googleTranslateService.GoogleTranslateV2DetectHandler())

	rtr.Methods("POST").Path("/google-translate/v3/translate").Handler(googleTranslateService.GoogleTranslateV3TranslateHandler())
	rtr.Methods("POST").Path("/google-translate/v3/detect").Handler(googleTranslateService.GoogleTranslateV3DetectHandler())

	rtr.Methods("POST").Path("/google-translate/translate").Handler(googleTranslateService.GoogleTranslateTranslateHandler())
	rtr.Methods("POST").Path("/google-translate/detect").Handler(googleTranslateService.GoogleTranslateDetectHandler())

	rtr.Methods("GET").Path("/google-translate/v3/glossaries").Handler(googleTranslateService.GoogleTranslateListGlossaryHandler())
	rtr.Methods("POST").Path("/google-translate/v3/glossaries").Handler(googleTranslateService.GoogleTranslateCreateGlossaryHandler())
	rtr.Methods("DELETE").Path("/google-translate/v3/glossaries").Handler(googleTranslateService.GoogleTranslateDeleteGlossaryHandler())

	registerMiddlewares(rtr)
	registerFallbackRoute(rtr)
}

func registerFallbackRoute(rtr *mux.Router) {
	rtr.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
	})
}

func registerMiddlewares(rtr *mux.Router) {
	rtr.Use(
		gziphandler.GzipHandler,
	)
}
