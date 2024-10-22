//services/user_services.go

package services

import (
    "errors"
    "auth-user-api/models"
    "auth-user-api/repository"
    "auth-user-api/utils"

    "golang.org/x/crypto/bcrypt"
)

type UserService interface {
    Register(username, email, password1, password2 string) error
    Update(id, username, email, password1, password2 string) error
    Delete(id string) error
    Authenticate(username, password string) error
    GetAllUsers() ([]*models.User, error)
    GetUserByID(id string) (*models.User, error)   // Tambahkan ini
}

type userService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userService{repo}
}

func (s *userService) Register(username, email, password1, password2 string) error {
    if password1 != password2 {
        return errors.New("password didn't match")
    }

    if err := utils.ValidatePassword(password1); err != nil {
        return err
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := &models.User{
        Username: username,
        Email:    email,
        Password: string(hashedPassword),
    }

    return s.repo.CreateUser(user)
}

func (s *userService) GetAllUsers() ([]*models.User, error) {
    users, err := s.repo.GetAllUsers()
    if err != nil {
        return nil, err
    }
    return users, nil
}

func (s *userService) Update(id, username, email, password1, password2 string) error {
    user, err := s.repo.GetUserByID(id)
    if err != nil {
        return err
    }

    if username != "" {
        user.Username = username
    }

    if email != "" {
        user.Email = email
    }

    if password1 != "" || password2 != "" {
        if password1 != password2 {
            return errors.New("password tidak cocok")
        }

        if err := utils.ValidatePassword(password1); err != nil {
            return err
        }

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
        if err != nil {
            return err
        }

        user.Password = string(hashedPassword)
    }

    return s.repo.UpdateUser(user)
}

func (s *userService) Delete(id string) error {
    return s.repo.DeleteUser(id)
}

func (s *userService) Authenticate(username, password string) error {
    user, err := s.repo.GetUserByUsername(username)
    if err != nil {
        if err.Error() == "record not found" {
            return errors.New("user not found")
        }
        return err
    }

    if user.DeletedAt.Valid {
        return errors.New("user not found")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return errors.New("Invalid username or password")
    }

    return nil
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
    return s.repo.GetUserByID(id)
}
