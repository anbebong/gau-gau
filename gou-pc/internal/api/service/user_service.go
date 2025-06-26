package service

import (
	"errors"
	"gou-pc/internal/api/model"
	"gou-pc/internal/api/repository"
	"time"
)

type UserService interface {
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	UpdatePassword(userID, newHash string) error
	DeleteUser(userID string) error
	AssignRole(userID, role string) error
	GetUserByID(userID string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	ListUsers() ([]model.User, error)
}

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) CreateUser(user *model.User) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.Username == user.Username {
			return errors.New("username already exists")
		}
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	users = append(users, *user)
	return s.repo.SaveAll(users)
}

func (s *userServiceImpl) UpdateUser(user *model.User) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	updated := false
	for i, u := range users {
		if u.ID == user.ID {
			user.UpdatedAt = time.Now()
			users[i] = *user
			updated = true
			break
		}
	}
	if !updated {
		return errors.New("user not found")
	}
	return s.repo.SaveAll(users)
}

func (s *userServiceImpl) UpdatePassword(userID, newHash string) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	updated := false
	for i, u := range users {
		if u.ID == userID {
			users[i].PasswordHash = newHash
			users[i].UpdatedAt = time.Now()
			updated = true
			break
		}
	}
	if !updated {
		return errors.New("user not found")
	}
	return s.repo.SaveAll(users)
}

func (s *userServiceImpl) DeleteUser(userID string) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	newUsers := users[:0]
	deleted := false
	for _, u := range users {
		if u.ID == userID {
			deleted = true
			continue
		}
		newUsers = append(newUsers, u)
	}
	if !deleted {
		return errors.New("user not found")
	}
	return s.repo.SaveAll(newUsers)
}

func (s *userServiceImpl) AssignRole(userID, role string) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	updated := false
	for i, u := range users {
		if u.ID == userID {
			users[i].Role = role
			users[i].UpdatedAt = time.Now()
			updated = true
			break
		}
	}
	if !updated {
		return errors.New("user not found")
	}
	return s.repo.SaveAll(users)
}

func (s *userServiceImpl) GetUserByID(userID string) (*model.User, error) {
	return s.repo.FindByID(userID)
}

func (s *userServiceImpl) GetUserByUsername(username string) (*model.User, error) {
	return s.repo.FindByUsername(username)
}

func (s *userServiceImpl) ListUsers() ([]model.User, error) {
	return s.repo.GetAll()
}
