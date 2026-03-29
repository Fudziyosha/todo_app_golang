package main

import (
	"context"
	"os"
	"os/signal"
	"web_todos/internal/config"
	"web_todos/internal/logger"
	"web_todos/internal/repository"
	"web_todos/internal/server"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.NewConfig()
	err := cfg.InitConfig()
	if err != nil {
		logrus.Error("main: failed init config ", err)
	}

	configPostgres := repository.NewPostgresConfig(
		viper.GetString("postgres.host"),
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.database"),
		viper.GetInt("postgres.port"))

	err = logger.InitLogger()
	if err != nil {
		logrus.Error("main: failed init config ", err)
	}

	repo := repository.NewRepository(configPostgres)
	err = repo.Connect(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	m := repository.NewMigration(configPostgres)
	err = m.UpMigration()
	if err != nil {
		logrus.Fatal("main: failed migrate db ", err)
	}

	fiber := server.NewServer(cfg)
	err = fiber.Server(repo)
	if err != nil {
		logrus.Fatal("main: failed up server ", err)
	}

}
