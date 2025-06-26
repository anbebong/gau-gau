package repository

import (
	"encoding/json"
	"errors"
	"gou-pc/internal/api/model"
	"gou-pc/internal/logutil"
	"os"
)

type UserRepository interface {
	GetAll() ([]model.User, error)
	SaveAll([]model.User) error
	FindByID(id string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
}

type fileUserRepository struct {
	file string
}

func NewFileUserRepository(file string) UserRepository {
	return &fileUserRepository{file: file}
}

func (r *fileUserRepository) GetAll() ([]model.User, error) {
	logutil.Debug("UserRepository.GetAll called")
	f, err := os.ReadFile(r.file)
	if err != nil {
		if os.IsNotExist(err) {
			logutil.Error("User DB file '%s' not found, returning empty list", r.file)
			return []model.User{}, nil
		}
		logutil.Debug("UserRepository.GetAll error: %v", err)
		return nil, err
	}
	var users []model.User
	if err := json.Unmarshal(f, &users); err != nil {
		logutil.Debug("UserRepository.GetAll unmarshal error: %v", err)
		return nil, err
	}
	logutil.Debug("UserRepository.GetAll loaded %d users", len(users))
	return users, nil
}

func (r *fileUserRepository) SaveAll(users []model.User) error {
	logutil.Debug("UserRepository.SaveAll called with %d users", len(users))
	b, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		logutil.Debug("UserRepository.SaveAll marshal error: %v", err)
		return err
	}
	return os.WriteFile(r.file, b, 0644)
}

func (r *fileUserRepository) FindByID(id string) (*model.User, error) {
	logutil.Debug("UserRepository.FindByID called with id=%s", id)
	users, err := r.GetAll()
	if err != nil {
		logutil.Debug("UserRepository.FindByID error: %v", err)
		return nil, err
	}
	for _, u := range users {
		logutil.Debug("UserRepository.FindByID checking user: id=%s", u.ID)
		if u.ID == id {
			logutil.Debug("UserRepository.FindByID found user: id=%s", u.ID)
			return &u, nil
		}
	}
	logutil.Debug("UserRepository.FindByID not found: id=%s", id)
	return nil, errors.New("user not found")
}

func (r *fileUserRepository) FindByUsername(username string) (*model.User, error) {
	logutil.Debug("UserRepository.FindByUsername called with username=%s", username)
	users, err := r.GetAll()
	if err != nil {
		logutil.Debug("UserRepository.FindByUsername error: %v", err)
		return nil, err
	}
	for _, u := range users {
		logutil.Debug("UserRepository.FindByUsername checking user: username=%s", u.Username)
		if u.Username == username {
			logutil.Debug("UserRepository.FindByUsername found user: username=%s", u.Username)
			return &u, nil
		}
	}
	logutil.Debug("UserRepository.FindByUsername not found: username=%s", username)
	return nil, errors.New("user not found")
}
