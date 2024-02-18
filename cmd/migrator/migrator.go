package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/synthao/notes/internal/config"
)

func main() {
	cnf, err := config.NewDBConfig()
	if err != nil {
		panic(err)
	}

	db := sqlx.MustConnect("mysql", cnf.GetMysqlDSN())
	driver, _ := mysql.WithInstance(db.DB, &mysql.Config{})

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
