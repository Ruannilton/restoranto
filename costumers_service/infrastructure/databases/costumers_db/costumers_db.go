package costumers_db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func ConnectPostgresDB() (*sql.DB, error) {

	dbURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		fmt.Println("Failed to connect to database")
		panic(err)
	}

	fmt.Println("Connected to database")

	return db, nil
}
