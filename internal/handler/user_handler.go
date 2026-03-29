package handler

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"web_todos/internal/entities"
	"web_todos/internal/repository"
	"web_todos/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/sirupsen/logrus"
)

const (
	sessionUserIDKey     = "user_id"
	sessionAuthenticated = "authenticated"
	FilePerm             = os.FileMode(0644)
	defaultAvatar        = "cat_avatar_default.jpg"
	defaultPathAvatar    = "./static/images/cat_avatar_default.jpg"
)

type UserHandler struct {
	repo     *repository.Repository
	validate *StructValidator
}

func NewUserHandler(repo *repository.Repository) *UserHandler {
	return &UserHandler{
		repo:     repo,
		validate: &StructValidator{Validator: validator.New()},
	}
}

func (u *UserHandler) GetRegistrationPage(c fiber.Ctx) error {
	return c.Render("register", nil)
}

func (u *UserHandler) UserRegistration(c fiber.Ctx) error {
	errorsField := make(map[string]string)

	user := new(entities.User)

	err := c.Bind().Form(user)
	if err != nil {
		logrus.Error("user_handler: failed parse form user ", err)
		return err
	}

	err = u.validate.Validate(user)
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			if e.Field() == "Password" {
				errorsField["Password"] = "Password minimum 8 character"
			}
			if e.Field() == "Email" {
				errorsField["Email"] = "Incorrect Email"
			}
		}
		return c.Render("register", fiber.Map{
			"Password": errorsField["Password"],
			"Email":    errorsField["Email"],
		})
	}

	hash, err := service.HashPassword(user.Password)
	if err != nil {
		logrus.Error("user_handler: failed hash password ", err)
		return err
	}

	userID, err := u.repo.User.CreateUser(c, user.Name, user.Surname, user.Email, hash, defaultAvatar, defaultPathAvatar)
	if err != nil {
		logrus.Error("user_handler: failed create user in handler ", err)
		return err
	}

	stringUserID := userID.String()

	err = setCookie(c, stringUserID)

	return c.Redirect().To("/")
}

func (u *UserHandler) UserLogin(c fiber.Ctx) error {
	user := new(entities.User)

	err := c.Bind().Form(user)
	if err != nil {
		logrus.Error("user_handler: failed get user form ", err)
		return err
	}

	err = u.validate.Validate(user)
	if err != nil {
		logrus.Error("user_handler: failed parse email address ", err)
		return c.Render("login", fiber.Map{
			"Email": "Incorrect email",
		})
	}

	_, hash, err := u.repo.User.GetUserIDAndPassword(c, user.Email)
	if err != nil {
		logrus.Error("user_handler: failed login ", err)
		return c.Render("login", fiber.Map{
			"Email": "Incorrect email",
		})
	}

	err = service.ValidatePassword(hash, user.Password)
	if err != nil {
		logrus.Error("user_handler: failed login ", err)
		return c.Render("login", fiber.Map{
			"Password": "Wrong password",
		})
	}

	userID, _, err := u.repo.User.GetUserIDAndPassword(c, user.Email)
	if err != nil {
		logrus.Error("user_handler: failed get user id ", err)
		return err
	}

	stringUserID := userID.String()

	err = setCookie(c, stringUserID)

	return c.Redirect().To("/")
}

func (u *UserHandler) Logout(c fiber.Ctx) error {
	sess := session.FromContext(c)

	err := sess.Destroy()
	if err != nil {
		logrus.Error("user_handler: failed destroy the session ", err)
		return err
	}
	return c.Redirect().To("/user/login")
}

func (u *UserHandler) GetUserLogin(c fiber.Ctx) error {
	return c.Render("login", nil)
}

func (u *UserHandler) ChangePassword(c fiber.Ctx) error {
	return c.Render("change_password", nil)
}

func (u *UserHandler) UpdateUserNameAndAvatar(c fiber.Ctx) error {
	userID, err := GetUserIDInSession(c)
	if err != nil {
		logrus.Error("user handler: failed parse uuid ", err)
		return err
	}

	name := c.FormValue("user_name")
	if name != "" {
		err = u.repo.User.UpdateUserName(c, name, userID)
		if err != nil {
			logrus.Error("user handler: failed update user name ", err)
			return err
		}
	}

	image, err := c.FormFile("avatar")
	if err == nil && image != nil {
		root, err := os.OpenRoot("./uploads")
		if err != nil {
			logrus.Error("user handler: failed open root dir ", err)
			return err
		}
		defer root.Close()

		imageFile, err := image.Open()
		bytes, err := io.ReadAll(imageFile)
		ext := filepath.Ext(image.Filename)

		imageName := fmt.Sprintf("avatar_%s%s", userID, ext)
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			logrus.Fatal("user handler: file is not image ")
		}

		if err := root.WriteFile("image_avatar/"+imageName, bytes, FilePerm); err != nil {
			logrus.Error("failed to write file: ", err)
			return err
		}

		file, err := root.Open("image_avatar/" + imageName)
		if err != nil {
			logrus.Error("unexpectedly succeeded in opening a file outside root ", err)

			err = root.Remove("image_avatar/" + imageName)
			if err != nil {
				logrus.Error("user handler: failed delete the avatar ", err)
				return err
			}

			return err
		}
		_ = file

		err = u.repo.User.UpdateImage(c, imageName, "./uploads/image_avatar/"+imageName, userID)
		if err != nil {
			logrus.Error("user handler: failed update image user ", err)
			return err
		}
	}

	return c.Redirect().To("/")
}

func (u *UserHandler) UpdateUserPass(c fiber.Ctx) error {
	userID, err := GetUserIDInSession(c)
	if err != nil {
		logrus.Error("user handler: failed parse uuid ", err)
		return err
	}

	currentPass := c.FormValue("current_password")

	user, err := u.repo.User.GetUser(c, userID)
	if err != nil {
		logrus.Error("user_handler: failed get user pass in db ", err)
		return err
	}

	err = service.ValidatePassword(user.Password, currentPass)
	if err != nil {
		logrus.Error("user_handler: pass is invalid ")
		return c.Render("change_password", fiber.Map{
			"Password": "Current password is not valid",
		})
	}

	val := validator.New()

	newPass := c.FormValue("user-password")
	err = val.Var(newPass, "required,min=8")
	if err != nil {
		return c.Render("change_password", fiber.Map{
			"Password": "Password minimum 8 character",
		})
	}

	confirmPass := c.FormValue("confirm_password")
	if newPass != confirmPass {
		return c.Render("change_password", fiber.Map{
			"UpdatePassword": "Wrong password",
		})
	}

	newHashPass, err := service.HashPassword(newPass)
	if err != nil {
		logrus.Error("user handler: failed hashed new pass ", err)
		return err
	}

	err = u.repo.User.UpdateUserPass(c, newHashPass, userID)

	return c.Redirect().To("/")
}

func setCookie(c fiber.Ctx, userID string) error {
	sess := session.FromContext(c)

	err := sess.Regenerate()
	if err != nil {
		logrus.Error("user_handler: failed regenerate session id ", err)
		return err
	}
	sess.Set(sessionUserIDKey, userID)
	sess.Set(sessionAuthenticated, true)

	return nil
}
