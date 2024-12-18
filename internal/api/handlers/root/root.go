package root

import (
	"net/http"

	"github.com/sabasm/go-server/internal/utils"
	"go.uber.org/zap"
)

const (
	methodNotAllowedMsg = "Method Not Allowed"
	welcomeMsg          = "Welcome to Go Server"
)

type RootResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type RootHandler struct {
	logger *zap.Logger
}

func New() (*RootHandler, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &RootHandler{
		logger: logger.Named("root-handler"),
	}, nil
}

func MustNew() *RootHandler {
	logger := zap.NewNop()
	return &RootHandler{
		logger: logger.Named("root-handler"),
	}
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	allowedMethods := map[string]bool{
		http.MethodGet:    true,
		http.MethodPost:   true,
		http.MethodPut:    true,
		http.MethodDelete: true,
		http.MethodPatch:  true,
	}

	if !allowedMethods[r.Method] {
		if err := utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, RootResponse{
			Error: methodNotAllowedMsg,
		}); err != nil {
			h.logger.Error("failed to write method not allowed response",
				zap.Error(err),
				zap.String("method", r.Method))
			return
		}
		return
	}

	if err := utils.WriteJSONResponse(w, http.StatusOK, RootResponse{
		Message: welcomeMsg,
	}); err != nil {
		h.logger.Error("failed to write success response",
			zap.Error(err),
			zap.String("method", r.Method))
	}
}
