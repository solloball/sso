package main

import (
    "flag"
    "fmt"
    "errors"

    "github.com/golang-migrate/migrate/v4"
    // driver for working with sqlite
    _ "github.com/golang-migrate/migrate/v4/database/sqlite3"
    // driver for working with file
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    var storagePath, migrationsPath, migrationsTable string

   flag.StringVar(&storagePath, "storage-path", "", "path to storage") 
   flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations") 
   flag.StringVar(
       &migrationsTable,
       "migraions_table",
       "migrations",
       "name of migrations table",
   ) 
   flag.Parse()

   mustValidateFlag(storagePath, migrationsPath)

   	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
        panic(err)
    }

    if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no change")
			return
		}

		panic(err)
	}
	fmt.Println("migrations finished successfully")
}

func mustValidateFlag(storagePath string, migrationsPath string) {
    if storagePath == "" {
        panic("storage-path is required")
    }
    if migrationsPath == "" {
        panic("migrations-path is required")
    }
}
