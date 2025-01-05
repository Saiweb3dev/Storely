// internal/service/user_service.go
package service

import (
    "context"
    "backend/internal/models"
    "backend/internal/repository"
    "backend/utils"
)

type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(user models.User) error {
    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        return err
    }
    user.Password = hashedPassword
    
    return s.repo.Create(context.Background(), &user)
}

func (s *UserService) AuthenticateUser(email, password string) (*models.User, error) {
    user, err := s.repo.FindByEmail(context.Background(), email)
    if err != nil {
        return nil, err
    }

    if err := utils.VerifyPassword(user.Password, password); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *UserService) GenerateToken(user *models.User) (string, error) {
    return utils.GenerateJWT(user.ID.Hex())
}