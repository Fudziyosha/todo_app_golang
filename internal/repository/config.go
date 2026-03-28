package repository

type PostgresConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

func NewPostgresConfig(host, user, password, databaseName string, port int) *PostgresConfig {
	return &PostgresConfig{
		Host:         host,
		Port:         port,
		User:         user,
		Password:     password,
		DatabaseName: databaseName,
	}
}
