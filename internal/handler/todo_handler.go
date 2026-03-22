package handler

import (
	"fmt"
	"time"
	"web_todos/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type TodoHandler struct {
	repo *repository.Repository
}

func NewTodoHandler(repo *repository.Repository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

func (t *TodoHandler) GetHome(c fiber.Ctx) error {
	userID, err := GetUserIdInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return err
	}

	listsSlice, err := t.repo.Todo.GetListsByID(c, userID)
	if err != nil {
		logrus.Error("todo handler: failed get list by user ", err)
		return err
	}

	for i, val := range listsSlice {
		sliceTasks, _ := t.repo.Todo.GetTodosByList(c, val.ID, false)
		listsSlice[i].Todos = append(val.Todos, sliceTasks...)
	}

	user, err := t.repo.User.GetUser(c, userID)
	if err != nil {
		logrus.Error("todo handler: failed get user ", err)
		return err
	}

	return c.Render("index", fiber.Map{
		"List":      listsSlice,
		"User":      user,
		"HomeLists": listsSlice,
	})
}

func (t *TodoHandler) CreateListInHomePage(c fiber.Ctx) error {
	userID, err := GetUserIdInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return err
	}

	todoDescription := c.FormValue("todo_new_text")
	if todoDescription != "" {
		return c.Redirect().To("/")
	}

	listDescription := c.FormValue("list_title")

	if listDescription != "" {
		err = t.repo.Todo.CreateList(c, listDescription, userID)
		if err != nil {
			logrus.Error("todo handler: failed insert repo ", err)
			return err
		}
	}

	return c.Redirect().To("/")
}

func (t *TodoHandler) GetTasksByUser(c fiber.Ctx) error {
	uuidPage := c.Params("id")
	listID, err := uuid.Parse(uuidPage)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid list ", err)
		return err
	}

	userID, err := GetUserIdInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return err
	}

	pageMode, taskStatus := t.checkFilters(c)

	listsSlice, err := t.repo.Todo.GetListsByID(c, userID)
	if err != nil {
		logrus.Error("todo handler: failed get lists and user ", err)
		return err
	}

	user, err := t.repo.User.GetUser(c, userID)

	todos, err := t.repo.Todo.GetTodosByList(c, listID, taskStatus)

	return c.Render("index", fiber.Map{
		"Todo":         todos,
		"List":         listsSlice,
		"ActiveListID": listID,
		"filters":      pageMode,
		"User":         user,
	})
}

func (t *TodoHandler) TaskHandler(c fiber.Ctx) error {
	uuidPage := c.Params("id")
	listID, err := uuid.Parse(uuidPage)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return err
	}

	userID, err := GetUserIdInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return err
	}

	listDescription := c.FormValue("list_title")

	if listDescription != "" {
		err := t.repo.Todo.CreateList(c, listDescription, userID)
		if err != nil {
			logrus.Error("handler: failed insert repo ", err)
			return err
		}
	}

	todoDescription := c.FormValue("todo_new_text")

	if todoDescription != "" {
		err := t.repo.Todo.InsertTodoByList(c, todoDescription, listID)
		if err != nil {
			logrus.Error("todo handler: failed insert todos ", err)
			return err
		}
	}

	updateTime := time.Now().UTC()

	var uuidTask uuid.UUID
	idTask := c.FormValue("todo_id")
	if idTask != "" {
		uuidTask, err = uuid.Parse(idTask)
		if err != nil {
			logrus.Error("todo_handler: failed parse uuid task ", err)
			return err
		}
	}

	action := c.FormValue("action")

	switch action {
	case "delete":
		err = t.repo.Todo.DeleteTodoById(c, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed delete task ", err)
			return err
		}
	case "update_todo":
		newDescription := c.FormValue("todo_text")

		err = t.repo.Todo.UpdateTodoDescriptionById(c, newDescription, updateTime, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed update todo ", err)
			return err
		}
	case "update_completed":
		err = t.repo.Todo.UpdateTodoStatusById(c, true, updateTime, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed update todo status completed ", err)
			return err
		}
	case "update_active":
		err = t.repo.Todo.UpdateTodoStatusById(c, false, updateTime, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed update todo status active ", err)
			return err
		}
	case "delete_list":
		err = t.repo.Todo.DeleteListById(c, listID)
		if err != nil {
			logrus.Error("todo_handler: failed delete list by id ", err)
		}
		return c.Redirect().To("/")
	}

	pageMode, _ := t.checkFilters(c)

	return c.Redirect().To(fmt.Sprintf("/list/%v/%s", listID, pageMode))
}

func (t *TodoHandler) CreateList(c fiber.Ctx) error {
	userId, err := GetUserIdInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return err
	}

	listDescription := c.FormValue("list_title")

	err = t.repo.Todo.CreateList(c, listDescription, userId)
	if err != nil {
		logrus.Error("todo handler: failed insert repo ", err)
		return err
	}

	lists, err := t.repo.Todo.GetListsByID(c, userId)
	if err != nil {
		logrus.Error("todo handler: failed get lists ", err)
		return err
	}

	return c.Render("index", fiber.Map{
		"List": lists,
	})
}

func (t *TodoHandler) checkFilters(c fiber.Ctx) (string, bool) {
	filters := c.Params("filters")

	var defaultBool = false
	if filters == "completed" {
		defaultBool = true
	}
	return filters, defaultBool
}
