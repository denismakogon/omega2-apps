package api

import "fmt"

type PostgresConfig struct {
	Host     string `json:"pg_host"`
	Port     string `json:"pg_port"`
	DB       string `json:"pg_db"`
	User     string `json:"pg_user"`
	Password string `json:"pg_pswd"`
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
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		pgsql.User, pgsql.Password,
		pgsql.Host, pgsql.Port, pgsql.DB)
}
