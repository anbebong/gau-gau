package service

import (
	"errors"
	"gou-pc/internal/api/model"
	"gou-pc/internal/api/repository"
	"gou-pc/internal/logutil"

	"github.com/google/uuid"
)

type UserService interface {
	UserGetAll() ([]model.User, error)
	UserGetByUsername(username string) (*model.User, error)
	UserGetByID(id string) (*model.User, error)
	UserCreate(user *model.User) error
	UserUpdate(user *model.User) error
	UserDeleteByUsername(username string) error
	UserDeleteByID(id string) error
}

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) UserGetAll() ([]model.User, error) {
	return s.repo.UserGetAll()
}

func (s *userServiceImpl) UserGetByUsername(username string) (*model.User, error) {
	return s.repo.UserFindByUsername(username)
}

func (s *userServiceImpl) UserGetByID(id string) (*model.User, error) {
	return s.repo.UserFindByID(id)
}

func (s *userServiceImpl) UserCreate(user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.FullName == "" {
		return errors.New("full_name is required")
	}
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	err := s.repo.UserCreate(user)
	if err != nil {
		logutil.APIDebug("UserService.Create: failed to create user %s: %v", user.Username, err)
		return err
	}
	logutil.APIDebug("UserService.Create: created user %s", user.Username)
	return nil
}

func (s *userServiceImpl) UserUpdate(user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.FullName == "" {
		return errors.New("full_name is required")
	}
	err := s.repo.UserUpdate(user)
	if err != nil {
		logutil.APIDebug("UserService.Update: failed to update user %s: %v", user.Username, err)
		return err
	}
	logutil.APIDebug("UserService.Update: updated user %s", user.Username)
	return nil
}

func (s *userServiceImpl) UserDeleteByUsername(username string) error {
	err := s.repo.UserDeleteByUsername(username)
	if err != nil {
		logutil.APIDebug("UserService.DeleteByUsername: failed to delete user %s: %v", username, err)
		return err
	}
	logutil.APIDebug("UserService.DeleteByUsername: deleted user %s", username)
	return nil
}

func (s *userServiceImpl) UserDeleteByID(id string) error {
	err := s.repo.UserDeleteByID(id)
	if err != nil {
		logutil.APIDebug("UserService.DeleteByID: failed to delete user id=%s: %v", id, err)
		return err
	}
	logutil.APIDebug("UserService.DeleteByID: deleted user id=%s", id)
	return nil
}
