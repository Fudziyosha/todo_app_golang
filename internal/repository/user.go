package repository

import (
	"context"
	"web_todos/internal/entities"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserPGRepository struct {
	repository *Repository
}

func NewUserRepository(repository *Repository) *UserPGRepository {
	return &UserPGRepository{repository: repository}
}

func (u *UserPGRepository) CreateUser(ctx context.Context, name, email, password, pathImage string) (uuid.UUID, error) {
	var userID uuid.UUID

	query := `INSERT INTO users(name, email, password, path_image) VALUES ($1 , $2, $3, $4 ) RETURNING id;`
	row := u.repository.database.QueryRow(ctx, query, name, email, password, pathImage)
	err := row.Scan(&userID)
	if err != nil {
		logrus.Error("user repository: failed create user ", err)
		return uuid.UUID{}, err
	}

	return userID, nil
}

func (u *UserPGRepository) GetUserIDAndPassword(ctx context.Context, email string) (userID uuid.UUID, hash string, err error) {
	query := `SELECT id, password FROM users WHERE email = ($1);`
	row := u.repository.database.QueryRow(ctx, query, email)
	err = row.Scan(&userID, &hash)
	if err != nil {
		logrus.Error("user repository: not found user in database uuid ", err)
		return userID, hash, err
	}

	return userID, hash, nil
}

func (u *UserPGRepository) GetUser(ctx context.Context, id uuid.UUID) (entities.User, error) {
	var user entities.User

	query := `SELECT id, name, email, password ,path_image FROM users WHERE id = $1;`
	row := u.repository.database.QueryRow(ctx, query, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.PathImage)
	if err != nil {
		logrus.Error("user repository: failed select all Todo ", err)
		return entities.User{}, err
	}

	return user, nil
}

func (u *UserPGRepository) UpdateImage(ctx context.Context, path string, id uuid.UUID) error {
	query := `UPDATE users SET path_image = $1 WHERE id = $2;`
	result, err := u.repository.database.Query(ctx, query, path, id)
	if err != nil {
		logrus.Error("user repository: failed create user ", err)
		return err
	}
	result.Close()

	return nil
}

func (u *UserPGRepository) UpdateUserName(ctx context.Context, name string, id uuid.UUID) error {
	query := `UPDATE users SET name = $1 WHERE id = $2;`

	result, err := u.repository.database.Query(ctx, query, name, id)
	if err != nil {
		logrus.Error("user repository: failed update user name ", err)
		return err
	}
	result.Close()

	return nil
}

func (u *UserPGRepository) UpdateUserPass(ctx context.Context, hash string, id uuid.UUID) error {
	query := `UPDATE users SET password = $1 WHERE id = $2;`

	result, err := u.repository.database.Query(ctx, query, hash, id)
	if err != nil {
		logrus.Error("user repository: failed update user pass ", err)
		return err
	}
	result.Close()

	return nil
}
