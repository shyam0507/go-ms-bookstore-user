package users_db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var Client *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "webastra"
	dbname   = "bookstore_users_db"
)

func init() {

	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	Client, err = sql.Open("postgres", dataSourceName)

	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}

	log.Print("Database connected successfully!")
}
