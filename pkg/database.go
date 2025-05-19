package pkg

import (
	"database/sql"
	"fmt"
	"log"
	"main/internal/config"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database interface {
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Beginx() (*sqlx.Tx, error)
	Begin() (*sql.Tx, error)
	Select(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	Exec(query string, args ...any) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
}

func NewPostgresDatabase(env config.Env) Database {
	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass
	dbName := env.DBName

	var dbURI string
	if dbUser == "" || dbPass == "" {
		dbURI = fmt.Sprintf("postgres://%s:%s/%s", dbHost, dbPort, dbName)
	} else {
		dbURI = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	}

	config, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		log.Fatal(err.Error())
	}

	nativeDB := stdlib.OpenDB(*config.ConnConfig)

	db := sqlx.NewDb(nativeDB, "pgx")

	err = db.Ping()
	if err != nil {
		log.Fatal("unable to connect to database: ", err.Error())
	}

	migration, err := os.ReadFile("scripts/migrate.sql")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = db.Exec(string(migration))
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
