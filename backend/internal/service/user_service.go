// internal/service/user_service.go
package service

import (
    "context"
    "time"
    "errors"
    "log"
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

var ErrAccountLocked = errors.New("account is locked")
var ErrInvalidCredentials = errors.New("Invalid Credentials. Account is locked")

func (s *UserService) HandleLoginAttempt(ctx context.Context, email, password, ipAddress string) (*models.User, error) {
    user, err := s.repo.FindByEmail(ctx, email)
    if err != nil {
        return nil, err
    }

    // Check lock status
    if user.IsLocked {
        if user.LockExpiresAt != nil && time.Now().Before(*user.LockExpiresAt) {
            return nil, ErrAccountLocked
        }
        // Reset lock if expired
        s.repo.UpdateLockStatus(ctx, user.ID, nil)
    }

    // Verify password
    if err := utils.VerifyPassword(user.Password, password); err != nil {
        // Handle failed attempt
        if user.FailedAttempts >= 4 {
            lockExpiry := time.Now().Add(15 * time.Minute)
            s.repo.UpdateLockStatus(ctx, user.ID, &lockExpiry)
        }
        s.repo.UpdateLoginStats(ctx, user.ID, ipAddress, false)
        return nil, ErrInvalidCredentials
    }

    // Update successful login stats
    err = s.repo.UpdateLoginStats(ctx, user.ID, ipAddress, true)
    if err != nil {
        log.Printf("Failed to update login stats: %v", err)
    }

    return user, nil
}