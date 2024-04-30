package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/icrowley/fake"
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

	studentID, err := createData(Student{Name: fake.FullName(), Age: 19})
	if err != nil {
		log.Fatal("Failed to insert data:", err)
	}
	fmt.Printf("ID of added student: %v\n", studentID)

	if err := readData(); err != nil {
		log.Fatal("Failed to read data:", err)
	}

	// student, err := updateData(Student{Name: "Deborah Garrett", Age: 30})
	// if err != nil {
	// 	log.Fatal("Failed to insert data:", err)
	// }
	// fmt.Printf("ID of added student: %v\n", student)

	// err = deleteData(Student{Age: 30, Name: "Deborah Garrett"})
	// if err != nil {
	// 	log.Fatal("Delete failed:", err)
	// }
	// fmt.Printf("Student deleted successfully")

	studentName := "Anthony Moore" // Example student name to filter by
	if err = innerJoin(Student{Name: studentName}); err != nil {
		log.Fatal(err)
	}

}

func initDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	studentTableSQL := `
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER
	);`
	if _, err := db.Exec(studentTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	courseTableSQL := `
	CREATE TABLE IF NOT EXISTS courses (
		course_id INTEGER PRIMARY KEY AUTOINCREMENT,
		course_name TEXT NOT NULL,
		student_id INTEGER,
		FOREIGN KEY(student_id) REFERENCES students(id)
	);
	
	INSERT INTO courses (course_name, student_id) VALUES ('Introduction to Programming', 7);
	INSERT INTO courses (course_name, student_id) VALUES ('Advanced Database Systems', 8);
	`
	if _, err := db.Exec(courseTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}

func createData(student Student) (int64, error) {
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

func readData() error {
	rows, err := db.Query("SELECT id, name, age FROM students")
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int
		if err := rows.Scan(&id, &name, &age); err != nil {
			return fmt.Errorf("scan error: %v", err)
		}
		fmt.Println(id, name, age)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows error: %v", err)
	}

	return nil
}

// func updateData(student Student) (Student, error) {

// 	result, err := db.Exec("UPDATE students SET age = ? WHERE name = ?", student.Age, student.Name)
// 	if err != nil {
// 		return Student{}, fmt.Errorf("updateData: %v", err)
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return Student{}, fmt.Errorf("updateData: getting rows affected: %v", err)
// 	}

// 	if rowsAffected == 0 {
// 		return Student{}, fmt.Errorf("updateData: no rows affected, student not found")
// 	}

// 	return student, nil
// }

// func deleteData(student Student) error {

// 	_, err := db.Exec("DELETE FROM students WHERE name = ?", student.Name)
// 	if err != nil {
// 		return fmt.Errorf("deleteStudentByName: %v", err)
// 	}
// 	return nil
// }

func innerJoin(student Student) error {
	query := `SELECT students.name, courses.course_name
	FROM students
	JOIN courses ON students.id = courses.student_id
	WHERE students.name = ?`

	rows, err := db.Query(query, student.Name)
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var courseName string
		if err := rows.Scan(&name, &courseName); err != nil {
			return fmt.Errorf("scan error: %v", err)
		}
		fmt.Println(name, courseName)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows error: %v", err)
	}

	return nil
}
