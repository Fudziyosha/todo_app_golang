package repository

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

type Migration struct {
	postgres *PostgresConfig
}

func NewMigration(postgres *PostgresConfig) *Migration {
	return &Migration{postgres: postgres}
}

func (m *Migration) UpMigration() error {
	migrations, err := migrate.New(
		"file://./migrations/postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", m.postgres.User, m.postgres.Password, m.postgres.Host, m.postgres.Port, m.postgres.DatabaseName))
	if err != nil {
		logrus.Error("repository: failed migrations db ", err)
	}
	err = migrations.Up()
	if err != nil {
		logrus.Error("repository: failed up ", err)
	}
	return nil
}
