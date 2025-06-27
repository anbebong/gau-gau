package repository

import (
	"database/sql"
	"errors"
	"gou-pc/internal/api/model"
	"gou-pc/internal/logutil"
	"time"
)

type UserRepository interface {
	UserGetAll() ([]model.User, error)
	UserFindByUsername(username string) (*model.User, error)
	UserFindByID(id string) (*model.User, error)
	UserCreate(user *model.User) error
	UserUpdate(user *model.User) error
	UserDeleteByUsername(username string) error
	UserDeleteByID(id string) error
}

type sqliteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepository{db: db}
}

func (r *sqliteUserRepository) UserGetAll() ([]model.User, error) {
	rows, err := r.db.Query(`SELECT id, username, password, email, full_name, role, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.FullName, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err == nil {
			users = append(users, u)
		}
	}
	return users, nil
}

func (r *sqliteUserRepository) UserFindByUsername(username string) (*model.User, error) {
	row := r.db.QueryRow(`SELECT id, username, password, email, full_name, role, created_at, updated_at FROM users WHERE username = ?`, username)
	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.FullName, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *sqliteUserRepository) UserFindByID(id string) (*model.User, error) {
	row := r.db.QueryRow(`SELECT id, username, password, email, full_name, role, created_at, updated_at FROM users WHERE id = ?`, id)
	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.FullName, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *sqliteUserRepository) UserCreate(user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.FullName == "" {
		return errors.New("full_name is required")
	}
	if user.Role == "" {
		return errors.New("role is required")
	}
	if user.ID == "" {
		return errors.New("id is required")
	}
	if user.CreatedAt == "" {
		return errors.New("created_at is required")
	}
	if user.UpdatedAt == "" {
		return errors.New("updated_at is required")
	}
	_, err := r.db.Exec(`INSERT INTO users (id, username, password, email, full_name, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Username, user.Password, user.Email, user.FullName, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		logutil.Debug("UserRepository.Create: failed to create user %s: %v", user.Username, err)
		return err
	}
	logutil.Debug("UserRepository.Create: created user %s", user.Username)
	return nil
}

func (r *sqliteUserRepository) UserUpdate(user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	// Build dynamic update query
	fields := []string{}
	args := []interface{}{}
	if user.Username != "" {
		fields = append(fields, "username=?")
		args = append(args, user.Username)
	}
	if user.Password != "" {
		fields = append(fields, "password=?")
		args = append(args, user.Password)
	}
	if user.Email != "" {
		fields = append(fields, "email=?")
		args = append(args, user.Email)
	}
	if user.FullName != "" {
		fields = append(fields, "full_name=?")
		args = append(args, user.FullName)
	}
	if user.Role != "" {
		fields = append(fields, "role=?")
		args = append(args, user.Role)
	}
	// Luôn cập nhật updated_at = time.Now()
	fields = append(fields, "updated_at=?")
	args = append(args, time.Now())
	if len(fields) == 0 {
		return errors.New("no fields to update")
	}
	args = append(args, user.ID)
	query := "UPDATE users SET " + joinFields(fields) + " WHERE id=?"
	_, err := r.db.Exec(query, args...)
	if err != nil {
		logutil.Debug("UserRepository.Update: failed to update user %s: %v", user.Username, err)
		return err
	}
	logutil.Debug("UserRepository.Update: updated user %s", user.Username)
	return nil
}

func joinFields(fields []string) string {
	if len(fields) == 0 {
		return ""
	}
	result := fields[0]
	for i := 1; i < len(fields); i++ {
		result += ", " + fields[i]
	}
	return result
}

func (r *sqliteUserRepository) UserDeleteByUsername(username string) error {
	res, err := r.db.Exec(`DELETE FROM users WHERE username=?`, username)
	n, _ := res.RowsAffected()
	if err != nil {
		logutil.Debug("UserRepository.DeleteByUsername: failed to delete user %s: %v", username, err)
		return err
	}
	if n == 0 {
		return errors.New("user not found")
	}
	logutil.Debug("UserRepository.DeleteByUsername: deleted user %s", username)
	return nil
}

func (r *sqliteUserRepository) UserDeleteByID(id string) error {
	res, err := r.db.Exec(`DELETE FROM users WHERE id=?`, id)
	n, _ := res.RowsAffected()
	if err != nil {
		logutil.Debug("UserRepository.DeleteByID: failed to delete user id=%s: %v", id, err)
		return err
	}
	if n == 0 {
		return errors.New("user not found")
	}
	logutil.Debug("UserRepository.DeleteByID: deleted user id=%s", id)
	return nil
}
