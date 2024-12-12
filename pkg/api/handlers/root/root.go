package root

import (
	"net/http"

	"github.com/sabasm/go-server/pkg/utils"
)

type RootHandler struct{}

func New() *RootHandler {
	return &RootHandler{}
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Welcome to Go Server"})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
