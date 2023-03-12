package googletranslatewrapper

import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	"github.com/weiyuan-lane/google-translate-api/internal/utils/errorhandlers"
	"golang.org/x/text/language"
)

type TranslateV2Wrapper struct {
	translateClient    *translate.Client
	translateV3Wrapper TranslateV3Wrapper
}

func NewTranslateV2Wrapper(translateClient *translate.Client) TranslateV2Wrapper {
	return TranslateV2Wrapper{
		translateClient: translateClient,
	}
}

func NewTranslateV2WrapperWithV3Wrapper(translateClient *translate.Client, translateV3Wrapper TranslateV3Wrapper) TranslateV2Wrapper {
	return TranslateV2Wrapper{
		translateClient:    translateClient,
		translateV3Wrapper: translateV3Wrapper,
	}
}

func (t TranslateV2Wrapper) TranslateText(ctx context.Context, text string, targetLocale language.Tag) (Translation, error) {
	googleTranslations, err := t.translateClient.Translate(
		ctx,
		[]string{text},
		targetLocale,
		&translate.Options{},
	)

	if err != nil {
		return Translation{}, errorhandlers.Wrap(
			errorhandlers.ErrGoogleTranslateV2EmptyTranslationResponse,
			fmt.Sprintf("Google translate returned error %s", err.Error()),
		)
	}

	if len(googleTranslations) == 0 {
		return Translation{}, errorhandlers.Wrap(
			errorhandlers.ErrGoogleTranslateV2EmptyTranslationResponse,
			"Google translate not returning any results",
		)
	}

	translation := t.makeTranslationResponse(googleTranslations[0], text, targetLocale)

	return translation, nil
}

func (t TranslateV2Wrapper) DetectionsFromText(ctx context.Context, text string) ([]Detection, error) {
	googleDetections, err := t.translateClient.DetectLanguage(
		ctx,
		[]string{text},
	)

	if err != nil {
		return []Detection{}, errorhandlers.Wrap(
			errorhandlers.ErrGoogleTranslateV2DetectErrResponse,
			fmt.Sprintf("Google translate detection returned error %s", err.Error()),
		)
	}

	if len(googleDetections) == 0 {
		return []Detection{}, errorhandlers.Wrap(
			errorhandlers.ErrGoogleTranslateV2EmptyDetectionResponse,
			"Google translate not returning any results",
		)
	}

	detections := t.makeDetectionsResponse(googleDetections[0])

	return detections, nil
}

func (t TranslateV2Wrapper) TranslateTextWithV3API(ctx context.Context, text string, targetLocale language.Tag, useV3API bool) (Translation, error) {
	if useV3API {
		v3TranslationResult, wrappedErr := t.translateV3Wrapper.TranslateText(ctx, text, targetLocale.String(), nil, nil)
		if wrappedErr != nil {
			return Translation{}, wrappedErr
		}

		detectedLangTag, err := language.Parse(v3TranslationResult.DetectedLang)
		if err != nil {
			return Translation{}, errorhandlers.Wrap(
				errorhandlers.ErrGoogleTranslateV2ConvertLangTagErrResponse,
				fmt.Sprintf("Google translate detected lang convert tag err %s", err.Error()),
			)
		}

		targetLangTag, err := language.Parse(v3TranslationResult.TargetLang)
		if err != nil {
			return Translation{}, errorhandlers.Wrap(
				errorhandlers.ErrGoogleTranslateV2ConvertLangTagErrResponse,
				fmt.Sprintf("Google translate target lang convert tag err %s", err.Error()),
			)
		}

		return Translation{
			TranslatedText: v3TranslationResult.TranslatedText,
			DetectedLang:   detectedLangTag,
			OriginalText:   v3TranslationResult.OriginalText,
			TargetLang:     targetLangTag,
		}, nil
	}

	return t.TranslateText(ctx, text, targetLocale)
}

func (t TranslateV2Wrapper) DetectionsFromTextWithV3API(ctx context.Context, text string, useV3API bool) ([]Detection, error) {
	if useV3API {
		v3DetectionResults, wrappedErr := t.translateV3Wrapper.DetectionsFromText(ctx, text)
		if wrappedErr != nil {
			return []Detection{}, wrappedErr
		}

		detectionResults := make([]Detection, len(v3DetectionResults))

		for i, v3DetectionResult := range v3DetectionResults {
			confidence := float64(v3DetectionResult.Confidence)
			isReliable := v3DetectionResult.Confidence > 0.8
			detectedLangTag, err := language.Parse(v3DetectionResult.Language)
			if err != nil {
				return []Detection{}, errorhandlers.Wrap(
					errorhandlers.ErrGoogleTranslateV2ConvertLangTagErrResponse,
					fmt.Sprintf("Google translate detect lang convert tag err %s", err.Error()),
				)
			}

			detectionResults[i] = Detection{
				Confidence: confidence,
				IsReliable: isReliable,
				Language:   detectedLangTag,
			}
		}

		return detectionResults, nil
	}

	return t.DetectionsFromText(ctx, text)
}

func (t TranslateV2Wrapper) makeTranslationResponse(googleTranslation translate.Translation, originalText string, targetLocale language.Tag) Translation {
	return Translation{
		TranslatedText: googleTranslation.Text,
		DetectedLang:   googleTranslation.Source,
		OriginalText:   originalText,
		TargetLang:     targetLocale,
	}
}

func (t TranslateV2Wrapper) makeDetectionsResponse(googleDetections []translate.Detection) []Detection {
	detections := make([]Detection, len(googleDetections))

	for i, googleDetection := range googleDetections {
		detections[i] = Detection{
			Confidence: googleDetection.Confidence,
			IsReliable: googleDetection.IsReliable,
			Language:   googleDetection.Language,
		}
	}

	return detections
}
