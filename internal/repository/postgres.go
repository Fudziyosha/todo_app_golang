package repository

import (
	"context"
	"fmt"

	pgxlogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/sirupsen/logrus"
)

type Postgres struct {
	postgres *PostgresConfig
}

func NewPostgres(postgres *PostgresConfig) *Postgres {
	return &Postgres{postgres: postgres}
}

func (p *Postgres) Connect(ctx context.Context) (*pgx.Conn, error) {
	connConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", p.postgres.User, p.postgres.Password, p.postgres.Host, p.postgres.Port, p.postgres.DatabaseName))
	if err != nil {
		logrus.Error("Unable to connect to database ", err)
	}

	connConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxlogrus.NewLogger(logrus.WithField("module", "db")),
		LogLevel: tracelog.LogLevel(logrus.ErrorLevel),
		Config:   tracelog.DefaultTraceLogConfig(),
	}

	connStr, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		logrus.Error("database: failed connect to database ", err)
	}

	return connStr, nil
}
