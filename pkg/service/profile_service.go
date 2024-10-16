package service

import (
	"music-saas/pkg/db"
	"music-saas/pkg/model"
)

type ProfileService struct {
	userRepo db.UserRepository
}

func NewProfileService(repo db.UserRepository) *ProfileService {
	return &ProfileService{userRepo: repo}
}

func (s *ProfileService) Profile(username string) (*model.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
