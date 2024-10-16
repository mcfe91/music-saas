package api

import (
	"encoding/json"
	"music-saas/internal/middleware"
	"music-saas/pkg/model"
	"music-saas/pkg/service"
	"net/http"
)

type ProfileAPI struct {
	profileService *service.ProfileService
}

func NewProfileAPI(profileService *service.ProfileService) *ProfileAPI {
	return &ProfileAPI{profileService: profileService}
}

func (a *ProfileAPI) Profile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.ContextUserKey).(*model.User)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	userData, err := a.profileService.Profile(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userData)
}
