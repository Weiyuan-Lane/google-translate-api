package httprequests

type GoogleTranslateTranslateRequestBody struct {
	Text         string `json:"text"`
	TargetLocale string `json:"target_locale"`
	UseV3API     bool   `json:"v3"`
	SourceLocale string `json:"source_locale"`
	Glossary     struct {
		ID string `json:"id"`
	} `json:"glossary"`
}

type GoogleTranslateDetectRequestBody struct {
	Text     string `json:"text"`
	UseV3API bool   `json:"v3"`
}

type GoogleTranslateCreateGlossaryBody struct {
	ID           string `json:"id"`
	GCSSource    string `json:"gcs_source"`
	SourceLocale string `json:"source_locale"`
	TargetLocale string `json:"target_locale"`
}

type GoogleTranslateDeleteGlossaryBody struct {
	ID string `json:"id"`
}
