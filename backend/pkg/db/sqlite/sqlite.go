package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

func openDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "file:pkg/db/sqlite/database.db")
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	migrations := &migrate.FileMigrationSource{
		Dir: "pkg/db/migrations",
	}

	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Applied %d migrations\n", n)
	fmt.Println("DB Operational")
	return db
}

var Db *sql.DB = openDatabase()
