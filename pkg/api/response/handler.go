package response

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	MsgInternalServerError = "internal server error"
)

func WriteJSONInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	zap.L().Error("internal error", zap.Error(err))
	WriteJSONResponse(w, r, http.StatusInternalServerError, ErrorResponse{Error: MsgInternalServerError})
}

func WriteJSONResponse(w http.ResponseWriter, _ *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if data == nil {
		_, err := io.WriteString(w, `{}`)
		if err != nil {
			zap.L().Error("can't write response", zap.Error(err))
		}
	} else {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			zap.L().Error("can't marshal/write response", zap.Error(err))
		}
	}
}

func HandleNotFoundError(w http.ResponseWriter, r *http.Request) {
	WriteJSONResponse(w, r, http.StatusNotFound, struct {
		Error string `json:"error"`
	}{http.StatusText(http.StatusNotFound)})
}

func HandleNoContentResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNoContent)
}
