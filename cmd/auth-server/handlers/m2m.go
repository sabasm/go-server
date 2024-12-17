package handlers

import (
	"net/http"
	"time"

	"github.com/sabasm/go-server/internal/utils"
)

type TokenResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

type M2MHandler struct {
	secretKey string
}

func NewM2MHandler(secretKey string) *M2MHandler {
	return &M2MHandler{secretKey: secretKey}
}

func (h *M2MHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := TokenResponse{
		Token:   h.secretKey,
		Expires: time.Now().Add(24 * time.Hour).Unix(),
	}

	if err := utils.WriteJSONResponse(w, http.StatusOK, response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
