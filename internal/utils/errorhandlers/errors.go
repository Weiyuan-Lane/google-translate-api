package errorhandlers

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload" // autoload env
)

var appName = os.Getenv("APP_NAME")

// All declared errors
var (
	ErrGoogleTranslateV2EmptyTranslationResponse     = fmt.Errorf("%s.%d", appName, 1)
	ErrGoogleTranslateV2TranslateErrResponse         = fmt.Errorf("%s.%d", appName, 2)
	ErrGoogleTranslateV2DetectErrResponse            = fmt.Errorf("%s.%d", appName, 3)
	ErrGoogleTranslateV2EmptyDetectionResponse       = fmt.Errorf("%s.%d", appName, 4)
	ErrTranslateEndpointMissingTextBodyParam         = fmt.Errorf("%s.%d", appName, 5)
	ErrTranslateEndpointMissingTargetLocaleBodyParam = fmt.Errorf("%s.%d", appName, 6)
	ErrDecodeJSONBodyFromRequestFailed               = fmt.Errorf("%s.%d", appName, 7)
	ErrEncodeJSONResponseFailed                      = fmt.Errorf("%s.%d", appName, 8)
	ErrDetectEndpointMissingTextBodyParam            = fmt.Errorf("%s.%d", appName, 9)
	ErrGoogleTranslateV3EmptyDetectionResponse       = fmt.Errorf("%s.%d", appName, 10)
	ErrGoogleTranslateV3DetectErrResponse            = fmt.Errorf("%s.%d", appName, 11)
	ErrGoogleTranslateV3EmptyTranslationResponse     = fmt.Errorf("%s.%d", appName, 12)
	ErrGoogleTranslateV3TranslateErrResponse         = fmt.Errorf("%s.%d", appName, 13)
	ErrGoogleTranslateV2ConvertLangTagErrResponse    = fmt.Errorf("%s.%d", appName, 14)
	ErrGoogleTranslateV3CreateGlossaryErrResponse    = fmt.Errorf("%s.%d", appName, 15)
	ErrGoogleTranslateV3ListGlossaryErrResponse      = fmt.Errorf("%s.%d", appName, 16)
	ErrGoogleTranslateV3DeleteGlossaryErrResponse    = fmt.Errorf("%s.%d", appName, 17)
	ErrGlossaryEndpointMissingIDBodyParam            = fmt.Errorf("%s.%d", appName, 18)
	ErrGlossaryEndpointMissingGCSSourceBodyParam     = fmt.Errorf("%s.%d", appName, 19)
)

// Categorized to slices
type errorPackage struct {
	HTTPStatusCode int
	Errors         []error
}

var (
	all400Errors = errorPackage{
		HTTPStatusCode: 400,
		Errors: []error{
			ErrDecodeJSONBodyFromRequestFailed,
		},
	}

	all404Errors = errorPackage{
		HTTPStatusCode: 404,
		Errors:         []error{},
	}

	all422Errors = errorPackage{
		HTTPStatusCode: 422,
		Errors: []error{
			ErrTranslateEndpointMissingTextBodyParam,
			ErrTranslateEndpointMissingTargetLocaleBodyParam,
			ErrDetectEndpointMissingTextBodyParam,
			ErrGlossaryEndpointMissingIDBodyParam,
			ErrGlossaryEndpointMissingGCSSourceBodyParam,
		},
	}

	all500Errors = errorPackage{
		HTTPStatusCode: 500,
		Errors: []error{
			ErrEncodeJSONResponseFailed,
		},
	}

	all502Errors = errorPackage{
		HTTPStatusCode: 502,
		Errors: []error{
			ErrGoogleTranslateV2TranslateErrResponse,
			ErrGoogleTranslateV2DetectErrResponse,
			ErrGoogleTranslateV2EmptyTranslationResponse,
			ErrGoogleTranslateV2EmptyDetectionResponse,
			ErrGoogleTranslateV3EmptyDetectionResponse,
			ErrGoogleTranslateV3DetectErrResponse,
			ErrGoogleTranslateV3EmptyTranslationResponse,
			ErrGoogleTranslateV3TranslateErrResponse,
			ErrGoogleTranslateV2ConvertLangTagErrResponse,
			ErrGoogleTranslateV3CreateGlossaryErrResponse,
			ErrGoogleTranslateV3ListGlossaryErrResponse,
			ErrGoogleTranslateV3DeleteGlossaryErrResponse,
		},
	}

	allErrorPackages = []errorPackage{
		all400Errors,
		all404Errors,
		all422Errors,
		all500Errors,
		all502Errors,
	}
)

const (
	defaultErrorStatusCode = 500
)
