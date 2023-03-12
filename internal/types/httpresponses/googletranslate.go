package httpresponses

type GoogleTranslateOriginalContent struct {
	Text           string `json:"text"`
	DetectedLocale string `json:"detected_locale"`
}

type GoogleTranslateTranslatedContent struct {
	Text   string `json:"text"`
	Locale string `json:"locale"`
}

type GoogleTranslateTranslatedResponse struct {
	OriginalContent           GoogleTranslateOriginalContent    `json:"original"`
	TranslatedContent         GoogleTranslateTranslatedContent  `json:"translated"`
	GlossaryTranslatedContent *GoogleTranslateTranslatedContent `json:"glossary_translated,omitempty"`
}

type GoogleTranslateDetectedLocale struct {
	Confidence float32 `json:"confidence"`
	Language   string  `json:"locale"`
}

type GoogleTranslateDetectedResponse struct {
	DetectedLocales []GoogleTranslateDetectedLocale `json:"results"`
}

type GoogleTranslateGlossary struct {
	ID        string `json:"id"`
	GCSSource string `json:"gcs_source"`
}

type GoogleTranslateListGlossariesResponse struct {
	Glossaries []GoogleTranslateGlossary `json:"glossaries"`
}
