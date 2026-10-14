package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"learn-platform/backend/internal/model"
	"learn-platform/backend/internal/repository"
)

type AuthService struct {
	authRepo *repository.AuthRepository
	userRepo *repository.UserRepository
	secret   string
}

func NewAuthService(
	authRepo *repository.AuthRepository,
	userRepo *repository.UserRepository,
	secret string,
) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		userRepo: userRepo,
		secret:   secret,
	}
}

func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, errors.New("missing required fields")
	}

	// Check if user exists
	if s.userRepo.ExistsByEmail(ctx, req.Email) {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Determine role
	role := req.Role
	if role == "" {
		role = "learner"
	}

	// Create user
	user := &model.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, refreshToken, expiresIn, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password required")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't leak if user exists or not (security best practice)
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if err := s.verifyPassword(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, refreshToken, expiresIn, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token required")
	}

	// Parse refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user id in token")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	accessToken, newRefreshToken, expiresIn, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user.ToResponse(),
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &model.Claims{
		ID:    claims["sub"].(string),
		Email: claims["email"].(string),
		Role:  claims["role"].(string),
	}, nil
}

func (s *AuthService) generateTokens(user *model.User) (string, string, int64, error) {
	accessExpire := time.Now().Add(24 * time.Hour)
	refreshExpire := time.Now().Add(7 * 24 * time.Hour)

	// Access token
	accessClaims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   accessExpire.Unix(),
		"iat":   time.Now().Unix(),
		"type": "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", 0, err
	}

	// Refresh token
	refreshClaims := jwt.MapClaims{
		"sub": user.ID,
		"exp": refreshExpire.Unix(),
		"iat": time.Now().Unix(),
		"type": "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", 0, err
	}

	return accessTokenString, refreshTokenString, int64(time.Until(accessExpire).Seconds()), nil
}

func (s *AuthService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *AuthService) verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
