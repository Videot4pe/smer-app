package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"path/filepath"
	"runtime"
)

//var embedMigrations embed.FS

func Migrate(dsn string) {

	var err error
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	//goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	_, b, _, _ := runtime.Caller(0)
	// TODO fix
	root := filepath.Join(filepath.Dir(b), "../../../..")
	migrations := filepath.Join(root, "migrations")

	if err := goose.Up(db, migrations); err != nil {
		panic(err)
	}
}
