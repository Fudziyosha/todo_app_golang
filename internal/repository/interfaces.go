package repository

import (
	"context"
	"time"
	"web_todos/internal/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	Todo TodoRepositoryHandler
	User UserRepositoryHandler
}

func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{
		Todo: NewTodoRepository(conn),
		User: NewUserRepository(conn),
	}
}

type TodoRepositoryHandler interface {
	CreateList(ctx context.Context, name string, userID uuid.UUID) error
	InsertTodoByList(ctx context.Context, description string, listID uuid.UUID) error
	GetListsById(ctx context.Context, userId uuid.UUID) ([]entities.List, error)
	GetTodosByList(ctx context.Context, id uuid.UUID, status bool) ([]entities.Todo, error)
	DeleteTodoById(ctx context.Context, todoId uuid.UUID) error
	UpdateTodoDescriptionById(ctx context.Context, newDescription string, timeUpdate time.Time, currentTodoId uuid.UUID) error
	UpdateTodoStatusById(ctx context.Context, status bool, timeUpdate time.Time, currentTodoId uuid.UUID) error
	DeleteListById(ctx context.Context, listID uuid.UUID) error
}

type UserRepositoryHandler interface {
	CreateUser(ctx context.Context, name, surname, email, password string) error
	UserAuth(ctx context.Context, email string) (hash string, err error)
	GetUserID(ctx context.Context, email string) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (entities.User, error)
	UpdateImage(ctx context.Context, name, path string, id uuid.UUID) error
	UpdateUserName(ctx context.Context, name string, id uuid.UUID) error
	UpdateUserPass(ctx context.Context, hash string, id uuid.UUID) error
}
