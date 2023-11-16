package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "your_username"
	password = "your_password"
	dbname   = "your_database_name"
)

var db *sql.DB

type Book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
}

func init() {
	var err error
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	row := db.QueryRow("SELECT * FROM books WHERE id = $1", id)

	var book Book
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	result, err := db.Exec("INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id", book.Title, book.Author)
	if err != nil {
		log.Fatal(err)
	}

	var createdID int
	err = result.Scan(&createdID)
	if err != nil {
		log.Fatal(err)
	}

	book.ID = createdID

	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	_, err := db.Exec("UPDATE books SET title=$1, author=$2 WHERE id=$3", book.Title, book.Author, id)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM books WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

