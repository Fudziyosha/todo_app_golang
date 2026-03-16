package repository

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func UpMigration() error {
	m, err := migrate.New(
		"file://./migrations/postgres",
		getDataBaseUrl())
	if err != nil {
		logrus.Error("repository: failed migrate db ", err)
	}
	err = m.Up()
	if err != nil {
		logrus.Error("repository: failed up ", err)
	}
	return nil
}
