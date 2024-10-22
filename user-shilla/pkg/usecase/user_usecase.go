package usecase

import (
	"errors"
	"project-golang-crud/domains"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	Repo domains.UserRepository
}

func NewUserUsecase(repo domains.UserRepository) domains.UserUsecase {
	return &userUsecase{Repo: repo}
}

func (u *userUsecase) Register(username, email, password string) (*domains.User, error) {
	var validationErrors []domains.ErrorDetail

	// Validasi username
	if err := validateUsername(username); err != nil {
		validationErrors = append(validationErrors, domains.ErrorDetail{
			Message: err.Error(),
			Parameter: "username",
		})
	}

	// Validasi email
	if err := validateEmail(email); err != nil {
		validationErrors = append(validationErrors, domains.ErrorDetail{
			Message: err.Error(),
			Parameter: "email",
		})
	}

	// Validasi password
	if err := validatePassword(password); err != nil {
		validationErrors = append(validationErrors, domains.ErrorDetail{
			Message: err.Error(),
			Parameter: "password",
		})
	}

	// Jika ada error validasi, return semua error
	if len(validationErrors) > 0 {
		return nil, errors.New(mergeValidationErrors(validationErrors))
	}

	// Cek apakah username sudah ada
	existingUser, err := u.Repo.GetByUsername(username)
	if err == nil && existingUser != nil {
		return nil, errors.New("duplicate username")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err // Handle hashing error
	}

	user := &domains.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	if err := u.Repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func mergeValidationErrors(errors []domains.ErrorDetail) string {
	var messages []string
	for _, err := range errors {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}


func mergeErrors(errors []string) string {
	return strings.Join(errors, "; ")
}

func (u *userUsecase) Update(id string, username, email, password string) error {
	user, err := u.Repo.GetByID(id)
	if err != nil {
		return domains.ErrUserNotFound
	}

	if user.DeletedAt != nil {
		return errors.New("User cannot be updated because it is marked as deleted")
	}

	var validationErrors []string

	// Validasi username
	if username == "" {
		return errors.New("Username is required")
	} else if username != user.Username {
		if err := validateUsername(username); err != nil {
			validationErrors = append(validationErrors, err.Error())
		} else if existingUser, _ := u.Repo.GetByUsername(username); existingUser != nil {
			validationErrors = append(validationErrors, "duplicate username")
		} else {
			user.Username = username // Update username
		}
	}

	// Validasi email
	if email != "" && email != user.Email {
		if err := validateEmail(email); err != nil {
			validationErrors = append(validationErrors, err.Error())
		} else {
			user.Email = email // Update email
		}
	}
	
	// Validasi password
	if password != "" {
		if err := validatePassword(password); err != nil {
			validationErrors = append(validationErrors, err.Error())
		} else {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return err // Handle hashing error
			}
			user.Password = string(hashedPassword) // Update password
		}
	}

	// Jika ada error validasi, return semua error
	if len(validationErrors) > 0 {
		return errors.New(mergeErrors(validationErrors))
	}

	return u.Repo.Update(user) // Lakukan pembaruan ke repositori
}

func (u *userUsecase) Delete(id string) (*domains.User, error) {
	user, err := u.Repo.GetByID(id)
	if err != nil {
		return nil, domains.ErrUserNotFound
	}

	if user.DeletedAt != nil {
		return nil, errors.New("User cannot be deleted because it is already marked as deleted")
	}

	if err := u.Repo.Delete(id); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) Validate(username, password string) (string, error) {
	user, err := u.Repo.GetByUsername(username)
	if err != nil {
		return "User not found", nil
	}

	// Cek apakah password valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "Invalid Password", nil // Password salah, tapi tidak melempar error
	}

	return "Valid Credentials", nil
}

func (u *userUsecase) GetByUsername(username string) (*domains.User, error) {
	user, err := u.Repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func validatePassword(password string) error {
	var (
		hasMinLen = len(password) >= 8 // Minimum length check
		hasNumber = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasUpper  = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'"<>,\./?\\|]`).MatchString(password)
	)
	if !hasMinLen || !hasNumber || !hasUpper || !hasSpecial{
		return errors.New("Password must be at least 8 characters long, contain an uppercase letter, a number, and a special character")
	}
	return nil
}

func validateUsername(username string) error {
	if match, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username); !match {
		return errors.New("Username can only contain letters, numbers, and underscores")
	}
	return nil
}

func validateEmail(email string) error {
	if match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email); !match {
		return errors.New("Invalid email format")
	}
	return nil
}

func (u *userUsecase) GetByID(id string) (*domains.User, error) {
	return u.Repo.GetByID(id)
}
