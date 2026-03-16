package repository

import (
	"context"
	"fmt"

	pgxlogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	connConfig, err := pgx.ParseConfig(getDataBaseUrl())
	if err != nil {
		logrus.Error("Unable to connect to database ", err)
	}

	connConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxlogrus.NewLogger(logrus.WithField("module", "db")),
		LogLevel: tracelog.LogLevel(logrus.InfoLevel),
		Config:   tracelog.DefaultTraceLogConfig(),
	}

	connStr, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		logrus.Error("database: failed connect to database ", err)
	}

	return connStr, nil
}

func getDataBaseUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.database"),
	)
}
