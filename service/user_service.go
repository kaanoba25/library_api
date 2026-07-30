package service

import (
	"errors"

	"github.com/kaanoba25/library_api/models"
	"github.com/kaanoba25/library_api/repository"
	"github.com/kaanoba25/library_api/utils"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(req models.RegisterRequest) (models.AuthResponse, error) {
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return models.AuthResponse{}, errors.New("all fields are required")
	}

	if req.Role != models.RoleAdmin && req.Role != models.RoleMember {
		req.Role = models.RoleMember // Default role
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.AuthResponse{}, err
	}

	user := models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	createdUser, err := s.repo.Create(user)
	if err != nil {
		return models.AuthResponse{}, errors.New("email already exists or database error")
	}

	token, err := utils.GenerateToken(createdUser.ID, string(createdUser.Role))
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{Token: token, User: createdUser}, nil
}

func (s *UserService) Login(req models.LoginRequest) (models.AuthResponse, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return models.AuthResponse{}, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return models.AuthResponse{}, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{Token: token, User: user}, nil
}