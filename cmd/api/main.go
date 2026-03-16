package main

import (
	"context"
	"os"
	"os/signal"
	"web_todos/internal/config"
	"web_todos/internal/logger"
	"web_todos/internal/repository"
	"web_todos/internal/serve"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := config.InitConfig()
	if err != nil {
		logrus.Error("main: failed init config ", err)
	}
	err = logger.InitLogger()
	if err != nil {
		logrus.Error("main: failed init config ", err)
	}

	conn, err := repository.Connect(ctx)
	if err != nil {
		logrus.Fatal("main: failed connect database ", err)
	}

	repo := repository.NewRepository(conn)

	err = repository.UpMigration()
	if err != nil {
		logrus.Fatal("main: failed migrate db ", err)
	}

	serve.Serve(repo)

}
