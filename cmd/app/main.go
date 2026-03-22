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

	err := config.InitConfig()
	if err != nil {
		logrus.Error("main: failed init config ", err)
	}
	err = logger.InitLogger()
	if err != nil {
		logrus.Error("main: failed init config ", err)
	}

	cfgPostgres := repository.NewPostgresConfig(viper.GetString("postgres.host"), viper.GetString("postgres.user"), viper.GetString("postgres.password"), viper.GetString("postgres.database"), viper.GetInt("postgres.port"))

	repo := repository.NewRepository(cfgPostgres)
	err = repo.Connect(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	m := repository.NewMigration(cfgPostgres)
	err = m.UpMigration()
	if err != nil {
		logrus.Fatal("main: failed migrate db ", err)
	}

	fiber := server.NewServer()
	fiber.Server(repo)

}
