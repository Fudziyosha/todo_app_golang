package repository

import (
	"context"
	"time"
	"web_todos/internal/entities"

	"github.com/google/uuid"
)

type TodoRepository interface {
	CreateList(ctx context.Context, name string, userID uuid.UUID) error
	InsertTodoByList(ctx context.Context, description string, listID uuid.UUID) error
	GetListsByID(ctx context.Context, userId uuid.UUID) ([]entities.List, error)
	GetTodosByListFilter(ctx context.Context, id uuid.UUID, filter string) ([]entities.Todo, error)
	DeleteTodoByID(ctx context.Context, todoId uuid.UUID) error
	UpdateTodoDescriptionByID(ctx context.Context, newDescription string, timeUpdate time.Time, currentTodoId uuid.UUID) error
	UpdateTodoStatusByID(ctx context.Context, status bool, timeUpdate time.Time, currentTodoId uuid.UUID) error
	DeleteListByID(ctx context.Context, listID []uuid.UUID) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, name, email, password, imageName, pathImage string) (uuid.UUID, error)
	GetUserIDAndPassword(ctx context.Context, email string) (userID uuid.UUID, hash string, err error)
	GetUser(ctx context.Context, id uuid.UUID) (entities.User, error)
	UpdateImage(ctx context.Context, name, path string, id uuid.UUID) error
	UpdateUserName(ctx context.Context, name string, id uuid.UUID) error
	UpdateUserPass(ctx context.Context, hash string, id uuid.UUID) error
}
