package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
	_ "modernc.org/sqlite"
)

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

var users []User

func createDb(db *sql.DB) {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(250),
		email VARCHAR(250) UNIQUE
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table 'users' created or already exists.")
}

func dbConn() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./db.db")
    if err != nil {
        return nil, err
    }
	fmt.Println("Connected to SQLite database!")
	return db, nil
}

func addUser(db *sql.DB, user User) (int, error) {
	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?,?)", user.Name, user.Email)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func getUser(db *sql.DB, id int) (User, error) {
	row := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)

	var user User
    err := row.Scan(&user.Id, &user.Name, &user.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return User{}, fmt.Errorf("user not found")
        }
        return User{}, err
    }

    return user, nil
}

func getAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
        return nil, err
    }
    defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan (&user.Id, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func main() {
    db, err := dbConn()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
	createDb(db)

    users = append(users, User{ Name: "John", Email: "1@2.3"})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Root")
	})
	
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := getAllUsers(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
	})
	
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/users/"):]
		id, err := strconv.Atoi(idStr)

		user, err := getUser(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
	})
	
	http.HandleFunc("/users/add", func(w http.ResponseWriter, r *http.Request) {
		if (r.Method != http.MethodPost) {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		var user User
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := addUser(db, user)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(id)
	})

	fmt.Println("Server is running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }

}
