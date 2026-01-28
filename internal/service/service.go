package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"test_task/internal/pkg/jwt"

	"test_task/internal/models"

	"test_task/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, id int64, req models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {

	existingUser, err := s.repo.GetByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this login already exists")
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FirstName: sql.NullString{String: req.FirstName, Valid: true},
		LastName:  sql.NullString{String: req.LastName, Valid: true},
		Login:     req.Login,
		Password:  hashedPassword,
		LastLogin: time.Now(),
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		user.Password = ""
	}
	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int64, req models.UpdateUserRequest) (*models.User, error) {

	user, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.FirstName != "" {
		user.FirstName = sql.NullString{String: req.FirstName, Valid: true}
	}
	if req.LastName != "" {
		user.LastName = sql.NullString{String: req.LastName, Valid: true}
	}
	if req.Login != "" {
		if req.Login != user.Login {
			existingUser, err := s.repo.GetByLogin(ctx, req.Login)
			if err != nil {
				return nil, err
			}
			if existingUser != nil && existingUser.ID != id {
				return nil, errors.New("login is already taken")
			}
			user.Login = req.Login
		}
	}
	if req.Password != "" {
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *userService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.repo.GetByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !checkPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	user.LastLogin = time.Now()
	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}
