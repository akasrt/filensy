package database

import (
	"database/sql"
	"log"

	"github.com/akasrt/filensy/internal/config/env"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	_ "github.com/go-sql-driver/mysql"
)

var db *sqlx.DB

func InitDB() {
	var err error

	dsn := env.GetEnv(env.DSN)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Println("unable to connect to mysql database")
		log.Fatal(err)
	}

	migrate(db.DB)
}

func GetDB() *sqlx.DB {
	return db
}

func migrate(db *sql.DB) error {
	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	return goose.Up(db, "../migrations")
}

func Close() {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
