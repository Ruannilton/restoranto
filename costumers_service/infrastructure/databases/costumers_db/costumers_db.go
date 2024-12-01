package costumers_db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func ConnectPostgresDB() (*sql.DB, error) {

	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	name := os.Getenv("PG_DB_NAME")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		fmt.Println("Failed to connect to database")
		panic(err)
	}

	fmt.Println("Connected to database")

	return db, nil
}
