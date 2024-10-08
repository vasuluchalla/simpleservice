package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	"encoding/json"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	API_PATH = "/apis/v1/books"
)

type Book struct {
	Id, Name, Isbn string
}

type databaseConfig struct {
	dbHost string
	dbUser string
	dbPass string
	dbName string
	dbPort string
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "food-app"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	l := databaseConfig{
		dbHost: dbHost,
		dbUser: dbUser,
		dbPass: dbPass,
		dbName: dbName,
		dbPort: dbPort,
	}
	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods(http.MethodGet)
	//r.HandleFunc(apiPath, l.getBooks).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)
}

func (l databaseConfig) getBooks(w http.ResponseWriter, r *http.Request) {
	db := l.openConnection()
	rows, err := db.Query("select * from books")
	if err != nil {
		panic(err)
	}

	books := []Book{}
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			panic(err)
		}

		abook := Book{Id: id, Name: name, Isbn: isbn}
		books = append(books, abook)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	}
}

func (l databaseConfig) openConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		l.dbHost, l.dbPort, l.dbUser, l.dbPass, l.dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("unable to connect to databse", err)
	}
	defer db.Close()
	return db
}

func (l databaseConfig) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("closing connection %s\n", err.Error())
	}
}
