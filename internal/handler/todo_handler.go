package handler

import (
	"fmt"
	"time"
	"web_todos/internal/entities"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetHome(c fiber.Ctx) error {
	b := h.CheckCookieAuthenticated(c)
	if b == false {
		return c.Redirect().To("/user/login")
	}

	userID, err := h.GetUserIdInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return err
	}

	listsSlice, user, err := h.selectListsUserTodos(c, userID, false)
	if err != nil {
		logrus.Error("todo handler: failed get lists and user ", err)
		return err
	}

	return c.Render("index", fiber.Map{
		"List":      listsSlice,
		"User":      user,
		"HomeLists": listsSlice,
	})
}

func (h *Handler) CreateListInHomePage(c fiber.Ctx) error {
	userID, err := h.GetUserIdInSession(c)
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
		err = h.repo.Todo.CreateList(c, listDescription, userID)
		if err != nil {
			logrus.Error("todo handler: failed insert repo ", err)
			return err
		}
	}

	return c.Redirect().To("/")
}

func (h *Handler) GetTasksByUser(c fiber.Ctx) error {
	uuidPage := c.Params("id")
	listID, err := uuid.Parse(uuidPage)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid list ", err)
		return err
	}

	userID, err := h.GetUserIdInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return err
	}

	pageMode, taskStatus := h.checkFilters(c)

	listsSlice, user, err := h.selectListsUserTodos(c, userID, taskStatus)
	if err != nil {
		logrus.Error("todo handler: failed select lists with todos ", err)
		return err
	}

	return c.Render("index", fiber.Map{
		"Todo":         listsSlice,
		"List":         listsSlice,
		"ActiveListID": listID,
		"filters":      pageMode,
		"User":         user,
	})
}

func (h *Handler) TaskHandler(c fiber.Ctx) error {
	uuidPage := c.Params("id")
	listID, err := uuid.Parse(uuidPage)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return err
	}

	userID, err := h.GetUserIdInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return err
	}

	listDescription := c.FormValue("list_title")

	if listDescription != "" {
		err := h.repo.Todo.CreateList(c, listDescription, userID)
		if err != nil {
			logrus.Error("handler: failed insert repo ", err)
			return err
		}
	}

	todoDescription := c.FormValue("todo_new_text")

	if todoDescription != "" {
		err := h.repo.Todo.InsertTodoByList(c, todoDescription, listID)
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
		err = h.repo.Todo.DeleteTodoById(c, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed delete task ", err)
			return err
		}
	case "update_todo":
		newDescription := c.FormValue("todo_text")

		err = h.repo.Todo.UpdateTodoDescriptionById(c, newDescription, updateTime, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed update todo ", err)
			return err
		}
	case "update_completed":
		err = h.repo.Todo.UpdateTodoStatusById(c, true, updateTime, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed update todo status completed ", err)
			return err
		}
	case "update_active":
		err = h.repo.Todo.UpdateTodoStatusById(c, false, updateTime, uuidTask)
		if err != nil {
			logrus.Error("todo_handler: failed update todo status active ", err)
			return err
		}
	case "delete_list":
		err = h.repo.Todo.DeleteListById(c, listID)
		if err != nil {
			logrus.Error("todo_handler: failed delete list by id ", err)
		}
		return c.Redirect().To("/")
	}

	pageMode, _ := h.checkFilters(c)

	return c.Redirect().To(fmt.Sprintf("/list/%v/%s", listID, pageMode))
}

func (h *Handler) CreateList(c fiber.Ctx) error {
	userId, err := h.GetUserIdInSession(c)
	if err != nil {
		logrus.Error("todo handler: failed parse uuid ", err)
		return err
	}

	listDescription := c.FormValue("list_title")

	err = h.repo.Todo.CreateList(c, listDescription, userId)
	if err != nil {
		logrus.Error("todo handler: failed insert repo ", err)
		return err
	}

	lists, err := h.repo.Todo.GetListsById(c, userId)
	if err != nil {
		logrus.Error("todo handler: failed get lists ", err)
		return err
	}

	return c.Render("index", fiber.Map{
		"List": lists,
	})
}

func (h *Handler) CheckCookieAuthenticated(c fiber.Ctx) bool {
	sess := session.FromContext(c)
	authenticated := sess.Get(sessionAuthenticated)

	if authenticated == nil {
		return false
	} else {
		return true
	}
}

func (h *Handler) checkFilters(c fiber.Ctx) (string, bool) {
	filters := c.Params("filters")

	var defaultBool = false
	if filters == "completed" {
		defaultBool = true
	}
	return filters, defaultBool
}

func (h *Handler) selectListsUserTodos(c fiber.Ctx, userID uuid.UUID, taskStatus bool) ([]entities.List, entities.User, error) {
	listsSlice, err := h.repo.Todo.GetListsById(c, userID)
	if err != nil {
		logrus.Error("todo handler: failed get list by user ", err)
		return []entities.List{}, entities.User{}, err
	}

	for i, val := range listsSlice {
		sliceTasks, _ := h.repo.Todo.GetTodosByList(c, val.ID, taskStatus)
		listsSlice[i].Todos = append(val.Todos, sliceTasks...)
	}

	user, err := h.repo.User.GetUser(c, userID)
	if err != nil {
		logrus.Error("todo handler: failed get user ", err)
		return []entities.List{}, entities.User{}, err
	}
	return listsSlice, user, nil
}
