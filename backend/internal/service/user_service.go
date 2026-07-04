package service

import (
	"context"
	"errors"
	"fmt"

	"learn-platform/backend/internal/model"
	"learn-platform/backend/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("invalid user id")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetAll(ctx context.Context, limit, offset int) ([]model.User, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Max limit for security
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *UserService) Update(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.User, error) {
	if id == "" {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user id")
	}
	return s.repo.Delete(ctx, id)
}
