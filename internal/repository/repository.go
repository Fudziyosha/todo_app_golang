package repository

import (
	"context"
	"fmt"

	pgxlogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	Todo     TodoRepository
	User     UserRepository
	database *pgx.Conn
	postgres *PostgresConfig
}

func NewRepository(postgres *PostgresConfig) *Repository {
	repository := &Repository{}

	repository.Todo = NewTodoRepository(repository)
	repository.User = NewUserRepository(repository)
	repository.postgres = postgres

	return repository
}
func (r *Repository) Connect(ctx context.Context) error {
	connConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", r.postgres.User, r.postgres.Password, r.postgres.Host, r.postgres.Port, r.postgres.DatabaseName))
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

	r.database = connStr

	return nil
}
