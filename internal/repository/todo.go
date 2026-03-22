package repository

import (
	"context"
	"time"
	"web_todos/internal/entities"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TodoPGRepository struct {
	repository *Repository
}

func NewTodoRepository(repository *Repository) *TodoPGRepository {
	return &TodoPGRepository{repository: repository}
}

func (t *TodoPGRepository) CreateList(ctx context.Context, name string, userID uuid.UUID) error {
	createdAt := time.Now().UTC()

	query := `INSERT INTO List(name,created_by, created_at) VALUES ($1 , $2, $3 );`
	resultInsert, err := t.repository.database.Query(ctx, query, name, userID, createdAt)
	if err != nil {
		logrus.Error("repository: failed scan insert List ", err)
		return err
	}

	resultInsert.Close()

	return nil
}

func (t *TodoPGRepository) InsertTodoByList(ctx context.Context, description string, listID uuid.UUID) error {
	createdAt := time.Now().UTC()

	query := `INSERT INTO Todo(description, status, created_at, created_in_list) VALUES ($1 , $2, $3, $4 );`
	resultInsert, err := t.repository.database.Query(ctx, query, description, false, createdAt, listID)
	if err != nil {
		logrus.Error("todo repository: failed scan insert List ", err)
		return err
	}
	resultInsert.Close()

	return nil
}

//func (t *TodoPGRepository) CreateListInHomePage(ctx context.Context, name string) error {
//	return nil
//}

func (t *TodoPGRepository) GetListsById(ctx context.Context, userId uuid.UUID) ([]entities.List, error) {
	query := `SELECT id, name, created_by, created_at, updated_at FROM List WHERE created_by = $1;`
	rows, err := t.repository.database.Query(ctx, query, userId)
	if err != nil {
		logrus.Error("todo repository: failed select all lists ", err)
		return nil, err
	}
	defer rows.Close()

	var list entities.List
	var lists []entities.List
	for rows.Next() {
		err = rows.Scan(&list.ID, &list.Name, &list.CreatedBy, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			logrus.Error("todo repository: failed rows scan all lists ", err)
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, err
}

func (t *TodoPGRepository) GetTodosByList(ctx context.Context, id uuid.UUID, status bool) ([]entities.Todo, error) {
	query := `SELECT id, description, status, updated_at, created_in_list FROM Todo WHERE created_in_list = $1 AND status = $2 ORDER BY created_at DESC ;`
	rows, err := t.repository.database.Query(ctx, query, id, status)
	if err != nil {
		logrus.Error("todo repository: failed select all Todo ", err)
		return nil, err
	}
	defer rows.Close()

	var todo entities.Todo
	var todos []entities.Todo
	for rows.Next() {
		err = rows.Scan(&todo.ID, &todo.Description, &todo.Status, &todo.UpdatedAt, &todo.CreatedInList)
		if err != nil {
			logrus.Error("todo repository: failed rows scan all Todo ", err)
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, err
}

func (t *TodoPGRepository) DeleteTodoById(ctx context.Context, todoId uuid.UUID) error {
	query := `DELETE FROM todo WHERE id = $1;`

	row, err := t.repository.database.Query(ctx, query, todoId)
	if err != nil {
		logrus.Error("todo repository: failed delete task ", err)
		return err
	}
	row.Close()

	return nil
}

func (t *TodoPGRepository) UpdateTodoDescriptionById(ctx context.Context, newDescription string, timeUpdate time.Time, currentTodoId uuid.UUID) error {
	query := `UPDATE todo SET description = $1, updated_at = $2 WHERE id = $3;`

	row, err := t.repository.database.Query(ctx, query, newDescription, timeUpdate, currentTodoId)
	if err != nil {
		logrus.Error("todo repository: failed delete task ", err)
		return err
	}
	row.Close()

	return nil
}

func (t *TodoPGRepository) UpdateTodoStatusById(ctx context.Context, status bool, timeUpdate time.Time, currentTodoId uuid.UUID) error {
	query := `UPDATE todo SET status = $1, updated_at = $2 WHERE id = $3;`

	row, err := t.repository.database.Query(ctx, query, status, timeUpdate, currentTodoId)
	if err != nil {
		logrus.Error("todo repository: failed delete task ", err)
		return err
	}
	row.Close()

	return nil
}

func (t *TodoPGRepository) DeleteListById(ctx context.Context, listID uuid.UUID) error {
	query := `DELETE FROM list WHERE id = $1;`

	row, err := t.repository.database.Query(ctx, query, listID)
	if err != nil {
		logrus.Error("todo repository: failed delete list ", err)
		return err
	}
	row.Close()

	return nil
}
