package googletranslate

import (
	"context"
	"net/http"

	"golang.org/x/text/language"

	"github.com/weiyuan-lane/google-translate-api/internal/services/googletranslatewrapper"
	"github.com/weiyuan-lane/google-translate-api/internal/types/httprequests"
	"github.com/weiyuan-lane/google-translate-api/internal/types/httpresponses"
	"github.com/weiyuan-lane/google-translate-api/internal/utils/errorhandlers"
	httputils "github.com/weiyuan-lane/google-translate-api/internal/utils/http"
	loggerutils "github.com/weiyuan-lane/google-translate-api/internal/utils/logger"
)

type GoogleTranslateService struct {
	Logger             *loggerutils.Logger
	TranslateV2Wrapper googletranslatewrapper.TranslateV2Wrapper
	TranslateV3Wrapper googletranslatewrapper.TranslateV3Wrapper
}

func (g GoogleTranslateService) GoogleTranslateV2TranslateHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		requestBody := httprequests.GoogleTranslateTranslateRequestBody{}

		// Decode http body logic
		wrappedErr := httputils.DecodeJSONBody(r, &requestBody)
		if wrappedErr != nil {
			g.Logger.Info(wrappedErr.Error())
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.Text == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrTranslateEndpointMissingTextBodyParam,
				"\"text\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.TargetLocale == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrTranslateEndpointMissingTextBodyParam,
				"\"target_locale\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		localeTag, err := language.Parse(requestBody.TargetLocale)
		if err != nil {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrTranslateEndpointMissingTextBodyParam,
				"Text locale in body is invalid",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Main service handler
		translation, wrappedErr := g.TranslateV2Wrapper.TranslateText(ctx, requestBody.Text, localeTag)
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Encoding for http response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		wrappedErr = httputils.EncodeJSONResponse(w, httpresponses.GoogleTranslateTranslatedResponse{
			OriginalContent: httpresponses.GoogleTranslateOriginalContent{
				Text:           translation.OriginalText,
				DetectedLocale: translation.DetectedLang.String(),
			},
			TranslatedContent: httpresponses.GoogleTranslateTranslatedContent{
				Text:   translation.TranslatedText,
				Locale: translation.TargetLang.String(),
			},
		})
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}
	}
}

func (g GoogleTranslateService) GoogleTranslateV2DetectHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		requestBody := httprequests.GoogleTranslateDetectRequestBody{}

		// Decode http body logic
		wrappedErr := httputils.DecodeJSONBody(r, &requestBody)
		if wrappedErr != nil {
			g.Logger.Info(wrappedErr.Error())
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.Text == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrDetectEndpointMissingTextBodyParam,
				"\"text\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Main service handler
		detections, wrappedErr := g.TranslateV2Wrapper.DetectionsFromText(ctx, requestBody.Text)
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Encoding for http response
		resDetections := make([]httpresponses.GoogleTranslateDetectedLocale, len(detections))
		for i, detection := range detections {
			resDetections[i] = httpresponses.GoogleTranslateDetectedLocale{
				Language:   detection.Language.String(),
				Confidence: float32(detection.Confidence),
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		wrappedErr = httputils.EncodeJSONResponse(w, httpresponses.GoogleTranslateDetectedResponse{
			DetectedLocales: resDetections,
		})
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}
	}
}

func (g GoogleTranslateService) GoogleTranslateV3TranslateHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		requestBody := httprequests.GoogleTranslateTranslateRequestBody{}

		// Decode http body logic
		wrappedErr := httputils.DecodeJSONBody(r, &requestBody)
		if wrappedErr != nil {
			g.Logger.Info(wrappedErr.Error())
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.Text == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrTranslateEndpointMissingTextBodyParam,
				"\"text\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.TargetLocale == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrTranslateEndpointMissingTextBodyParam,
				"\"target_locale\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		var glossaryIDPtr *string
		if requestBody.Glossary.ID != "" {
			glossaryID := requestBody.Glossary.ID
			glossaryIDPtr = &glossaryID
		} else {
			glossaryIDPtr = nil
		}

		var sourceLangPtr *string
		if requestBody.SourceLocale != "" {
			sourceLang := requestBody.SourceLocale
			sourceLangPtr = &sourceLang
		} else {
			sourceLangPtr = nil
		}

		// Main service handler
		translation, wrappedErr := g.TranslateV3Wrapper.TranslateText(ctx, requestBody.Text, requestBody.TargetLocale, sourceLangPtr, glossaryIDPtr)
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Encoding for http response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		translatedResponse := httpresponses.GoogleTranslateTranslatedResponse{
			OriginalContent: httpresponses.GoogleTranslateOriginalContent{
				Text:           translation.OriginalText,
				DetectedLocale: translation.DetectedLang,
			},
			TranslatedContent: httpresponses.GoogleTranslateTranslatedContent{
				Text:   translation.TranslatedText,
				Locale: translation.TargetLang,
			},
		}
		if translation.GlossaryTranslatedText != "" {
			translatedResponse.GlossaryTranslatedContent = &httpresponses.GoogleTranslateTranslatedContent{
				Text:   translation.GlossaryTranslatedText,
				Locale: translation.TargetLang,
			}
		}

		wrappedErr = httputils.EncodeJSONResponse(w, translatedResponse)
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}
	}
}

func (g GoogleTranslateService) GoogleTranslateV3DetectHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		requestBody := httprequests.GoogleTranslateDetectRequestBody{}

		// Decode http body logic
		wrappedErr := httputils.DecodeJSONBody(r, &requestBody)
		if wrappedErr != nil {
			g.Logger.Info(wrappedErr.Error())
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.Text == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrDetectEndpointMissingTextBodyParam,
				"\"text\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Main service handler
		detections, wrappedErr := g.TranslateV3Wrapper.DetectionsFromText(ctx, requestBody.Text)
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Encoding for http response
		resDetections := make([]httpresponses.GoogleTranslateDetectedLocale, len(detections))
		for i, detection := range detections {
			resDetections[i] = httpresponses.GoogleTranslateDetectedLocale{
				Language:   detection.Language,
				Confidence: detection.Confidence,
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		wrappedErr = httputils.EncodeJSONResponse(w, httpresponses.GoogleTranslateDetectedResponse{
			DetectedLocales: resDetections,
		})
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}
	}
}

func (g GoogleTranslateService) GoogleTranslateDetectHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		requestBody := httprequests.GoogleTranslateDetectRequestBody{}

		// Decode http body logic
		wrappedErr := httputils.DecodeJSONBody(r, &requestBody)
		if wrappedErr != nil {
			g.Logger.Info(wrappedErr.Error())
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.Text == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrDetectEndpointMissingTextBodyParam,
				"\"text\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Main service handler
		detections, wrappedErr := g.TranslateV2Wrapper.DetectionsFromTextWithV3API(
			ctx,
			requestBody.Text,
			requestBody.UseV3API,
		)
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Encoding for http response
		resDetections := make([]httpresponses.GoogleTranslateDetectedLocale, len(detections))
		for i, detection := range detections {
			resDetections[i] = httpresponses.GoogleTranslateDetectedLocale{
				Language:   detection.Language.String(),
				Confidence: float32(detection.Confidence),
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		wrappedErr = httputils.EncodeJSONResponse(w, httpresponses.GoogleTranslateDetectedResponse{
			DetectedLocales: resDetections,
		})
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}
	}
}

func (g GoogleTranslateService) GoogleTranslateTranslateHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		requestBody := httprequests.GoogleTranslateTranslateRequestBody{}

		// Decode http body logic
		wrappedErr := httputils.DecodeJSONBody(r, &requestBody)
		if wrappedErr != nil {
			g.Logger.Info(wrappedErr.Error())
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.Text == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrTranslateEndpointMissingTextBodyParam,
				"\"text\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.TargetLocale == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrTranslateEndpointMissingTextBodyParam,
				"\"target_locale\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		localeTag, err := language.Parse(requestBody.TargetLocale)
		if err != nil {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrTranslateEndpointMissingTextBodyParam,
				"Text locale in body is invalid",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Main service handler
		translation, wrappedErr := g.TranslateV2Wrapper.TranslateTextWithV3API(
			ctx,
			requestBody.Text,
			localeTag,
			requestBody.UseV3API,
		)
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Encoding for http response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		wrappedErr = httputils.EncodeJSONResponse(w, httpresponses.GoogleTranslateTranslatedResponse{
			OriginalContent: httpresponses.GoogleTranslateOriginalContent{
				Text:           translation.OriginalText,
				DetectedLocale: translation.DetectedLang.String(),
			},
			TranslatedContent: httpresponses.GoogleTranslateTranslatedContent{
				Text:   translation.TranslatedText,
				Locale: translation.TargetLang.String(),
			},
		})
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}
	}
}

func (g GoogleTranslateService) GoogleTranslateCreateGlossaryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		requestBody := httprequests.GoogleTranslateCreateGlossaryBody{}

		// Decode http body logic
		wrappedErr := httputils.DecodeJSONBody(r, &requestBody)
		if wrappedErr != nil {
			g.Logger.Info(wrappedErr.Error())
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.ID == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrGlossaryEndpointMissingIDBodyParam,
				"\"id\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.GCSSource == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrGlossaryEndpointMissingGCSSourceBodyParam,
				"\"gcs_source\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Main service handler
		wrappedErr = g.TranslateV3Wrapper.CreateGlossary(
			ctx,
			requestBody.ID,
			requestBody.GCSSource,
			requestBody.SourceLocale,
			requestBody.TargetLocale,
		)
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Encoding for http response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("{}"))
	}
}

func (g GoogleTranslateService) GoogleTranslateDeleteGlossaryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		requestBody := httprequests.GoogleTranslateDeleteGlossaryBody{}

		// Decode http body logic
		wrappedErr := httputils.DecodeJSONBody(r, &requestBody)
		if wrappedErr != nil {
			g.Logger.Info(wrappedErr.Error())
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		if requestBody.ID == "" {
			wrappedErr := errorhandlers.Wrap(
				errorhandlers.ErrGlossaryEndpointMissingIDBodyParam,
				"\"id\" field in body is empty",
			)
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Main service handler
		wrappedErr = g.TranslateV3Wrapper.DeleteGlossary(
			ctx,
			requestBody.ID,
		)
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		// Encoding for http response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}

func (g GoogleTranslateService) GoogleTranslateListGlossaryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// Main service handler
		glossaries, wrappedErr := g.TranslateV3Wrapper.ListGlossaries(ctx)
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}

		glossaryResults := make([]httpresponses.GoogleTranslateGlossary, len(glossaries))
		for i, glossary := range glossaries {
			glossaryResults[i] = httpresponses.GoogleTranslateGlossary{
				ID:        glossary.ID,
				GCSSource: glossary.GCSSource,
			}
		}

		// Encoding for http response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		wrappedErr = httputils.EncodeJSONResponse(w, httpresponses.GoogleTranslateListGlossariesResponse{
			Glossaries: glossaryResults,
		})
		if wrappedErr != nil {
			errorhandlers.HandleHTTPError(g.Logger, wrappedErr, w)
			return
		}
	}
}
