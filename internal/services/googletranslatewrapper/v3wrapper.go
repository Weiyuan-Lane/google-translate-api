package googletranslatewrapper

import (
	"context"
	"fmt"

	translate "cloud.google.com/go/translate/apiv3"
	translatepb "cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/weiyuan-lane/google-translate-api/internal/utils/errorhandlers"
	"google.golang.org/api/iterator"
)

type TranslateV3Wrapper struct {
	translateClient *translate.TranslationClient
	projectKey      string
}

func NewTranslateV3Wrapper(translateClient *translate.TranslationClient, projectKey string) TranslateV3Wrapper {
	return TranslateV3Wrapper{
		translateClient: translateClient,
		projectKey:      projectKey,
	}
}

func (t TranslateV3Wrapper) TranslateText(ctx context.Context, text, targetLocale string, sourceLocale, glossaryID *string) (TranslationV3, error) {
	req := &translatepb.TranslateTextRequest{
		Parent:             t.projectKey, // Required
		MimeType:           "text/plain",
		Contents:           []string{text},
		TargetLanguageCode: targetLocale,
	}

	// Use glossary if indicated
	if glossaryID != nil {
		req.GlossaryConfig = &translatepb.TranslateTextGlossaryConfig{
			Glossary: *glossaryID,
		}
	}

	if sourceLocale != nil {
		req.SourceLanguageCode = *sourceLocale
	}

	googleTranslationResponse, err := t.translateClient.TranslateText(
		ctx,
		req,
	)

	if err != nil {
		return TranslationV3{}, errorhandlers.Wrap(
			errorhandlers.ErrGoogleTranslateV3EmptyTranslationResponse,
			fmt.Sprintf("Google translate returned error %s", err.Error()),
		)
	}

	if googleTranslationResponse == nil || len(googleTranslationResponse.Translations) == 0 {
		return TranslationV3{}, errorhandlers.Wrap(
			errorhandlers.ErrGoogleTranslateV3EmptyTranslationResponse,
			"Google translate not returning any results",
		)
	}

	translation, wrappedErr := t.makeTranslationResponse(
		googleTranslationResponse.Translations,
		googleTranslationResponse.GlossaryTranslations,
		text,
		targetLocale,
	)
	if wrappedErr != nil {
		return TranslationV3{}, wrappedErr
	}

	return translation, nil
}

func (t TranslateV3Wrapper) DetectionsFromText(ctx context.Context, text string) ([]DetectionV3, error) {
	req := &translatepb.DetectLanguageRequest{
		Parent:   t.projectKey, // Required
		MimeType: "text/plain",
		Source: &translatepb.DetectLanguageRequest_Content{
			Content: text,
		},
	}

	googleDetectionResponse, err := t.translateClient.DetectLanguage(
		ctx,
		req,
	)

	if err != nil {
		return []DetectionV3{}, errorhandlers.Wrap(
			errorhandlers.ErrGoogleTranslateV3DetectErrResponse,
			fmt.Sprintf("Google translate detection returned error %s", err.Error()),
		)
	}

	if googleDetectionResponse == nil {
		return []DetectionV3{}, errorhandlers.Wrap(
			errorhandlers.ErrGoogleTranslateV3EmptyDetectionResponse,
			"Google translate not returning any results",
		)
	}

	detections := t.makeDetectionsResponse(*googleDetectionResponse)

	return detections, nil
}

func (t TranslateV3Wrapper) CreateGlossary(ctx context.Context, id, gcsSource, sourceLocale, targetLocale string) error {
	glossary := &translatepb.Glossary{
		Name:        id,
		DisplayName: id,
		InputConfig: &translatepb.GlossaryInputConfig{
			Source: &translatepb.GlossaryInputConfig_GcsSource{
				GcsSource: &translatepb.GcsSource{
					InputUri: gcsSource,
				},
			},
		},
		Languages: &translatepb.Glossary_LanguagePair{
			LanguagePair: &translatepb.Glossary_LanguageCodePair{
				SourceLanguageCode: sourceLocale,
				TargetLanguageCode: targetLocale,
			},
		},
	}
	req := &translatepb.CreateGlossaryRequest{
		Parent:   t.projectKey, // Required
		Glossary: glossary,
	}

	_, err := t.translateClient.CreateGlossary(
		ctx,
		req,
	)
	if err != nil {
		return errorhandlers.Wrap(
			errorhandlers.ErrGoogleTranslateV3CreateGlossaryErrResponse,
			fmt.Sprintf("Google translate create glossary returning error: %s", err.Error()),
		)
	}

	return nil
}

func (t TranslateV3Wrapper) ListGlossaries(ctx context.Context) ([]GlossariesV3, error) {
	req := &translatepb.ListGlossariesRequest{
		Parent: t.projectKey,
	}

	glossaries := t.translateClient.ListGlossaries(
		ctx,
		req,
	)

	results := []GlossariesV3{}
	for currItem, err := glossaries.Next(); err != iterator.Done; currItem, err = glossaries.Next() {
		if err != nil {
			return []GlossariesV3{}, errorhandlers.Wrap(
				errorhandlers.ErrGoogleTranslateV3ListGlossaryErrResponse,
				fmt.Sprintf("Google translate list glossary returning error: %s", err.Error()),
			)
		}

		results = append(results, GlossariesV3{
			ID:        currItem.Name,
			GCSSource: currItem.InputConfig.GetGcsSource().GetInputUri(),
		})
	}

	return results, nil
}

func (t TranslateV3Wrapper) DeleteGlossary(ctx context.Context, id string) error {
	req := &translatepb.DeleteGlossaryRequest{
		Name: id,
	}

	_, err := t.translateClient.DeleteGlossary(
		ctx,
		req,
	)
	if err != nil {
		return errorhandlers.Wrap(
			errorhandlers.ErrGoogleTranslateV3DeleteGlossaryErrResponse,
			fmt.Sprintf("Google translate delete glossary returning error: %s", err.Error()),
		)
	}

	return nil
}

func (t TranslateV3Wrapper) makeTranslationResponse(
	googleTranslations []*translatepb.Translation,
	googleGlossaryTranslations []*translatepb.Translation,
	originalText string,
	targetLocale string,
) (TranslationV3, error) {

	if len(googleTranslations) == 0 {
		return TranslationV3{}, errorhandlers.Wrap(
			errorhandlers.ErrGoogleTranslateV3EmptyTranslationResponse,
			"Google translate not returning any results",
		)
	}

	translationResponse := TranslationV3{
		TranslatedText: googleTranslations[0].TranslatedText,
		DetectedLang:   googleTranslations[0].DetectedLanguageCode,
		OriginalText:   originalText,
		TargetLang:     targetLocale,
	}

	if len(googleGlossaryTranslations) > 0 {
		translationResponse.GlossaryTranslatedText = googleGlossaryTranslations[0].TranslatedText
	}

	return translationResponse, nil
}

func (t TranslateV3Wrapper) makeDetectionsResponse(googleDetectionResponse translatepb.DetectLanguageResponse) []DetectionV3 {
	detections := make([]DetectionV3, len(googleDetectionResponse.Languages))

	for i, googleDetection := range googleDetectionResponse.Languages {
		detections[i] = DetectionV3{
			Confidence: googleDetection.Confidence,
			Language:   googleDetection.LanguageCode,
		}
	}

	return detections
}
