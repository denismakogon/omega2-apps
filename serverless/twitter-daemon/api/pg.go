package api

import "fmt"

type PostgresConfig struct {
	Host     string `json:"postgres_host"`
	Port     string `json:"postgres_port"`
	DB       string `json:"postgres_db"`
	User     string `json:"postgres_user"`
	Password string `json:"postgres_password"`
}

func (pgsql *PostgresConfig) FromFile() error {
	envVar := "POSTGRES_APPLICATION_CONFIG"
	return StructFromFile(pgsql, envVar)
}

func (pgsql *PostgresConfig) FromEnv() error {
	return StructFromEnv(pgsql)
}

func (pgsql *PostgresConfig) DNS() string {
	return fmt.Sprintf(
		"dbname=%s user=%s password=%s port=%s host=%s sslmode=disable",
		pgsql.DB, pgsql.User,
		pgsql.Password, pgsql.Port, pgsql.Host)
}
