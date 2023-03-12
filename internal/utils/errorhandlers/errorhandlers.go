package errorhandlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	errorwrapper "github.com/pkg/errors"
	"github.com/weiyuan-lane/google-translate-api/internal/types/httpresponses"
	loggerutils "github.com/weiyuan-lane/google-translate-api/internal/utils/logger"
)

func HandleHTTPError(logger *loggerutils.Logger, stackErr error, w http.ResponseWriter) {

	var messageErr, baseErr error
	errResponse := httpresponses.ErrorResponse{}
	statusCode := defaultErrorStatusCode

	// Retrieve status code
	for _, errorPackage := range allErrorPackages {
		if errorIs(stackErr, errorPackage.Errors) {
			statusCode = errorPackage.HTTPStatusCode
			break
		}
	}

	// Get error contents
	messageErr = errors.Unwrap(stackErr)
	if messageErr != nil {
		baseErr = errors.Unwrap(messageErr)
	}

	if baseErr != nil {
		id := baseErr.Error()
		message := messageErr.Error()
		stack := fmt.Sprintf("%+v", stackErr)
		errResponse = httpresponses.ErrorResponse{
			ErrorCode: httpresponses.ErrorCode{
				ID:          id,
				Description: message,
			},
		}

		// Log error results
		logger.Info(message,
			map[string]string{
				"error_code_id": id,
				"error_stack":   stack,
			},
		)
	} else {
		logger.Info("Error does not contain any message",
			map[string]string{
				"error_code_id": "",
				"error_stack":   stackErr.Error(),
			},
		)
	}

	// Render response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errResponse)
}

func Wrap(err error, msg string) error {
	return errorwrapper.Wrap(err, msg)
}

func errorIs(err error, errorEntities []error) bool {
	result := false

	for _, errorEntity := range errorEntities {
		if errors.Is(err, errorEntity) {
			result = true
			break
		}
	}

	return result
}
