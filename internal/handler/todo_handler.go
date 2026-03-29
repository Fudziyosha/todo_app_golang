package handler

import (
	"fmt"
	"strings"
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
	userID, err := GetUserIDInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return c.SendStatus(500)
	}

	listsSlice, err := t.repo.Todo.GetListsByID(c, userID)
	if err != nil {
		logrus.Error("todo handler: failed get list by user ", err)
		return c.SendStatus(500)
	}

	for i, val := range listsSlice {
		sliceTasks, _ := t.repo.Todo.GetTodosByListFilter(c, val.ID, "active")
		listsSlice[i].Todos = append(val.Todos, sliceTasks...)
	}

	user, err := t.repo.User.GetUser(c, userID)
	if err != nil {
		logrus.Error("todo handler: failed get user ", err)
		return c.SendStatus(500)
	}

	return c.Render("index", fiber.Map{
		"List":      listsSlice,
		"User":      user,
		"HomeLists": listsSlice,
	})
}

func (t *TodoHandler) CreateListInHomePage(c fiber.Ctx) error {
	userID, err := GetUserIDInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return c.SendStatus(500)
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
			return c.SendStatus(400)
		}
	}

	if c.FormValue("action") != "" {
		listsID, err := getSliceLists(c)

		err = t.repo.Todo.DeleteListByID(c, listsID)
		if err != nil {
			logrus.Error("todo_handler: failed delete list by id ", err)
			return c.SendStatus(500)
		}
	}

	return c.Redirect().To("/")
}

func (t *TodoHandler) GetTasksByUser(c fiber.Ctx) error {
	uuidPage := c.Params("id")
	listID, err := uuid.Parse(uuidPage)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid list ", err)
		return c.SendStatus(500)
	}

	userID, err := GetUserIDInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return c.SendStatus(500)
	}

	query := c.Query("filters")

	todos, err := t.repo.Todo.GetTodosByListFilter(c, listID, query)
	if err != nil {
		logrus.Error("todo handler: failed get todos by list filter ", err)
		return c.SendStatus(500)
	}

	listsSlice, err := t.repo.Todo.GetListsByID(c, userID)
	if err != nil {
		logrus.Error("todo handler: failed get lists and user ", err)
		return c.SendStatus(500)
	}

	user, err := t.repo.User.GetUser(c, userID)
	*user.PathImage = strings.TrimPrefix(*user.PathImage, ".")

	return c.Render("index", fiber.Map{
		"Todo":         todos,
		"List":         listsSlice,
		"ActiveListID": listID,
		"filters":      query,
		"User":         user,
	})
}

func (t *TodoHandler) TaskHandler(c fiber.Ctx) error {
	uuidPage := c.Params("id")
	listID, err := uuid.Parse(uuidPage)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return c.SendStatus(500)
	}

	userID, err := GetUserIDInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return c.SendStatus(500)
	}

	listDescription := c.FormValue("list_title")

	if listDescription != "" {
		err := t.repo.Todo.CreateList(c, listDescription, userID)
		if err != nil {
			logrus.Error("handler: failed insert repo ", err)
			return c.SendStatus(400)
		}
	}

	todoDescription := c.FormValue("todo_new_text")

	if todoDescription != "" {
		err := t.repo.Todo.InsertTodoByList(c, todoDescription, listID)
		if err != nil {
			logrus.Error("todo handler: failed insert todos ", err)
			return c.SendStatus(400)
		}
	}

	updateTime := time.Now().UTC()

	var uuidTask uuid.UUID
	idTask := c.FormValue("todo_id")
	if idTask != "" {
		uuidTask, err = uuid.Parse(idTask)
		if err != nil {
			logrus.Error("todo_handler: failed parse uuid task ", err)
			return c.SendStatus(500)
		}
	}

	action := c.FormValue("action")

	switch action {
	case "delete":
		err = t.repo.Todo.DeleteTodoByID(c, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed delete task ", err)
			return c.SendStatus(500)
		}
	case "update_todo":
		newDescription := c.FormValue("todo_text")

		err = t.repo.Todo.UpdateTodoDescriptionByID(c, newDescription, updateTime, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed update todo ", err)
			return c.SendStatus(400)
		}
	case "update_completed":
		err = t.repo.Todo.UpdateTodoStatusByID(c, true, updateTime, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed update todo status completed ", err)
			return c.SendStatus(500)
		}
	case "update_active":
		err = t.repo.Todo.UpdateTodoStatusByID(c, false, updateTime, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed update todo status active ", err)
			return c.SendStatus(500)
		}
	case "delete_list":
		listsID, err := getSliceLists(c)

		err = t.repo.Todo.DeleteListByID(c, listsID)
		if err != nil {
			logrus.Error("todo_handler: failed delete list by id ", err)
			return c.SendStatus(500)
		}
		return c.Redirect().To("/")
	}

	taskBool := c.Query("filters")

	return c.Redirect().To(fmt.Sprintf("/list/%v?filters=%s", listID, taskBool))
}

func (t *TodoHandler) CreateList(c fiber.Ctx) error {
	userId, err := GetUserIDInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return c.SendStatus(500)
	}

	listDescription := c.FormValue("list_title")

	err = t.repo.Todo.CreateList(c, listDescription, userId)
	if err != nil {
		logrus.Error("todo handler: failed insert repo ", err)
		return c.SendStatus(400)
	}

	lists, err := t.repo.Todo.GetListsByID(c, userId)
	if err != nil {
		logrus.Error("todo handler: failed get lists ", err)
		return c.SendStatus(500)
	}

	return c.Render("index", fiber.Map{
		"List": lists,
	})
}

func getSliceLists(c fiber.Ctx) (listsID []uuid.UUID, err error) {

	lists, err := c.MultipartForm()
	for _, stringID := range lists.Value["list_id"] {
		listUUID, _ := uuid.Parse(stringID)
		listsID = append(listsID, listUUID)
	}

	return listsID, nil
}
