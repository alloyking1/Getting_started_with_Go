package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Student struct {
	Name string
	Age  int
}

var db *sql.DB

func main() {
	var err error
	db, err = initDB("./school.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database.")

	studentID, err := insertData(Student{Name: "student one", Age: 19})
	if err != nil {
		log.Fatal("Failed to insert data:", err)
	}

	fmt.Printf("ID of added student: %v\n", studentID)
}

func initDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER
	);`
	if _, err := db.Exec(sqlStmt); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}

func insertData(student Student) (int64, error) {
	result, err := db.Exec("INSERT INTO students (name, age) VALUES (?, ?)", student.Name, student.Age)
	if err != nil {
		return 0, fmt.Errorf("insertData: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insertData: %v", err)
	}
	return id, nil
}
