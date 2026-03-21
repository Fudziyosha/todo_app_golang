package repository

import (
	"context"
	"web_todos/internal/entities"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	repository *Repository
}

func NewUserRepository(repository *Repository) *UserRepository {
	return &UserRepository{repository: repository}
}

func (u *UserRepository) CreateUser(ctx context.Context, name, surname, email, password string) error {
	query := `INSERT INTO users(name, surname, email, password) VALUES ($1 , $2, $3, $4 );`
	result, err := u.repository.database.Query(ctx, query, name, surname, email, password)
	if err != nil {
		logrus.Error("user repository: failed create user ", err)
		return err
	}
	result.Close()

	return nil
}

func (u *UserRepository) UserAuth(ctx context.Context, email string) (hash string, err error) {
	query := `SELECT password FROM users WHERE email = ($1);`
	row := u.repository.database.QueryRow(ctx, query, email)
	err = row.Scan(&hash)
	if err != nil {
		logrus.Error("user repository: not found user in database ", err)
		return "", err
	}

	return hash, nil
}

func (u *UserRepository) GetUserID(ctx context.Context, email string) (uuid.UUID, error) {
	var userID uuid.UUID

	query := `SELECT id FROM users WHERE email = ($1);`
	row := u.repository.database.QueryRow(ctx, query, email)
	err := row.Scan(&userID)
	if err != nil {
		logrus.Error("user repository: not found user in database uuid ", err)
		return uuid.UUID{}, err
	}

	return userID, nil
}

func (u *UserRepository) GetUser(ctx context.Context, id uuid.UUID) (entities.User, error) {
	query := `SELECT id, name, surname, email, password ,image_name,path_image FROM users WHERE id = $1;`
	rows, err := u.repository.database.Query(ctx, query, id)
	if err != nil {
		logrus.Error("user repository: failed select all Todo ", err)
		return entities.User{}, err
	}
	defer rows.Close()

	var user entities.User
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.ImageName, &user.PathImage)
		if err != nil {
			logrus.Error("user repository: failed rows scan all Todo ", err)
			return entities.User{}, err
		}
	}
	return user, nil
}

func (u *UserRepository) UpdateImage(ctx context.Context, name, path string, id uuid.UUID) error {
	query := `UPDATE users SET image_name = $1, path_image = $2 WHERE id = $3;`
	result, err := u.repository.database.Query(ctx, query, name, path, id)
	if err != nil {
		logrus.Error("user repository: failed create user ", err)
		return err
	}
	result.Close()

	return nil
}

func (u *UserRepository) UpdateUserName(ctx context.Context, name string, id uuid.UUID) error {
	query := `UPDATE users SET name = $1 WHERE id = $2;`

	result, err := u.repository.database.Query(ctx, query, name, id)
	if err != nil {
		logrus.Error("user repository: failed update user name ", err)
		return err
	}
	result.Close()

	return nil
}

func (u *UserRepository) UpdateUserPass(ctx context.Context, hash string, id uuid.UUID) error {
	query := `UPDATE users SET password = $1 WHERE id = $2;`

	result, err := u.repository.database.Query(ctx, query, hash, id)
	if err != nil {
		logrus.Error("user repository: failed update user pass ", err)
		return err
	}
	result.Close()

	return nil
}
