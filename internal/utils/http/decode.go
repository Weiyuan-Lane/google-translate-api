package http

import (
	"encoding/json"
	"net/http"

	"github.com/weiyuan-lane/google-translate-api/internal/utils/errorhandlers"
)

func DecodeJSONBody(r *http.Request, ptr interface{}) error {
	err := json.NewDecoder(r.Body).Decode(ptr)
	if err != nil {
		wrappedErr := errorhandlers.Wrap(
			errorhandlers.ErrDecodeJSONBodyFromRequestFailed,
			err.Error(),
		)

		return wrappedErr
	}

	return nil
}
