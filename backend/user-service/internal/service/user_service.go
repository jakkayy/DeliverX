package service

import (
	"errors"

	"github.com/deliverx/user-service/internal/model"
	"github.com/deliverx/user-service/internal/repository"
	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	Name         string  `json:"name"`
	Email        *string `json:"email"`
	ProfileImage *string `json:"profile_image"`
	FcmToken     *string `json:"fcm_token"`
}

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfile(userID uuid.UUID) (*model.User, error) {
	return s.repo.FindByID(userID)
}

func (s *UserService) UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != nil {
		user.Email = req.Email
	}
	if req.ProfileImage != nil {
		user.ProfileImage = req.ProfileImage
	}
	if req.FcmToken != nil {
		user.FcmToken = req.FcmToken
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetDrivers() ([]model.User, error) {
	return s.repo.FindDrivers()
}
