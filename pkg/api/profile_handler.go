package api

import (
	"encoding/json"
	"music-saas/internal/middleware"
	"music-saas/pkg/model"
	"music-saas/pkg/service"
	"net/http"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.ContextUserKey).(*model.User)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	userData, err := h.profileService.Profile(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userData)
}
