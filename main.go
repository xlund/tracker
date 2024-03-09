package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xlund/tracker/repo"
)

var dbFile = "tracker.db"

func main() {
	port := flag.String("port", "3000", "Port to run the server on")
	flag.Parse()
	mux := http.NewServeMux()
	LoadHandlers(mux)

	http.ListenAndServe(":"+*port, mux)
}

func LoadHandlers(router *http.ServeMux) {
	us := LoadServices()
	us.Repo.Migrate()

	router.HandleFunc("/", Index())
	router.HandleFunc("GET /users/new", NewUser(us))
	router.HandleFunc("DELETE /users/{id}", DeleteUser(us))
	router.HandleFunc("GET /users", ListUsers(us))
}

func LoadServices() *repo.UserService {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	userRepo := repo.NewSQLiteUserRepo(db)
	userService := repo.NewUserService(userRepo)
	return userService
}

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the home page!")
	}
}

func NewUser(us *repo.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := us.Repo.Create(repo.User{
			Username: "xlund",
			Email:    "xlund@gmail.com",
			Name:     "Erik Xlund",
		})
		if err != nil {
			fmt.Fprintln(w, "Error creating user")
		}
		fmt.Fprintln(w, "User created")
	}
}

func ListUsers(us *repo.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := us.Repo.FindAll()
		if err != nil {
			fmt.Fprintln(w, "Error fetching users")
			return
		}

		for _, user := range users {
			fmt.Fprintf(w, "Username: %s, Email: %s, Name: %s\n", user.Username, user.Email, user.Name)
		}
	}
}

func DeleteUser(us *repo.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.PathValue("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			fmt.Fprintln(w, "Error parsing id")
			return
		}
		err = us.Repo.Delete(id)
		if err != nil {
			fmt.Fprintln(w, "Error deleting user")
			return
		}
		fmt.Fprintln(w, "User deleted")
	}
}
