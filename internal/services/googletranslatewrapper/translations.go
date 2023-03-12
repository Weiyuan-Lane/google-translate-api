package googletranslatewrapper

import (
	"golang.org/x/text/language"
)

type Translation struct {
	TranslatedText string
	OriginalText   string
	DetectedLang   language.Tag
	TargetLang     language.Tag
}

type Detection struct {
	Confidence float64
	IsReliable bool
	Language   language.Tag
}

type TranslationV3 struct {
	TranslatedText         string
	OriginalText           string
	DetectedLang           string
	TargetLang             string
	GlossaryTranslatedText string
}

type DetectionV3 struct {
	Confidence float32
	Language   string
}

type GlossariesV3 struct {
	ID        string
	GCSSource string
}
