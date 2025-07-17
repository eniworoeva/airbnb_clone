package service

import (
	"airbnb-clone/internal/logger"
	"errors"

	"airbnb-clone/internal/config"
	"airbnb-clone/internal/models"
	"airbnb-clone/internal/repository"
	"airbnb-clone/internal/utils"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo   repository.UserRepository
	jwtManager *utils.JWTManager
}

func NewUserService(userRepo repository.UserRepository, jwtConfig config.JWTConfig) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtManager: utils.NewJWTManager(jwtConfig),
	}
}

type LoginResponse struct {
	User         *models.UserResponse `json:"user"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
}

func (s *UserService) Register(req *models.UserCreateRequest) (*models.UserResponse, error) {
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("failed to check existing user: %v", err)
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Errorf("failed to hash password: %v", err)
		return nil, err
	}

	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      req.Role,
		IsActive:  true,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		logger.Errorf("failed to create user: %v", err)
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *UserService) Login(req *models.UserLoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		logger.Errorf("failed to get user by email: %v", err)
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		logger.Errorf("failed to generate access token: %v", err)
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		logger.Errorf("failed to generate refresh token: %v", err)
		return nil, err
	}

	return &LoginResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) ValidateToken(tokenString string) (*utils.Claims, error) {
	return s.jwtManager.ValidateToken(tokenString)
}

func (s *UserService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		logger.Errorf("failed to get user: %v", err)
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		logger.Errorf("failed to generate access token: %v", err)
		return nil, err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		logger.Errorf("failed to generate refresh token: %v", err)
		return nil, err

	}

	return &LoginResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
