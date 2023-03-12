package http

import (
	"encoding/json"
	"net/http"

	"github.com/weiyuan-lane/google-translate-api/internal/utils/errorhandlers"
)

func EncodeJSONResponse(w http.ResponseWriter, res interface{}) error {
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		wrappedErr := errorhandlers.Wrap(
			errorhandlers.ErrEncodeJSONResponseFailed,
			err.Error(),
		)

		return wrappedErr
	}

	return nil
}
