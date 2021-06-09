package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/prianka01/iitk-coin/user/register"

	"github.com/prianka01/iitk-coin/user/login"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)
func main() {
    r := mux.NewRouter()
	r.HandleFunc("/signup", register.RegisterHandler).
		Methods("POST")
	r.HandleFunc("/login", login.LoginHandler).
		Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
    database, _ := sql.Open("sqlite3", "./userData.db")
 	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, Name TEXT, Rollno INTEGER, Password TEXT, Token TEXT)")
    statement.Exec()
	// addUser("Priyanka",190649,"random");
	//  rows, _ := database.Query("SELECT id, name, rollno FROM user")
    // var id int
    // var name string
    // var rollno int
    // for rows.Next() {
    //     rows.Scan(&id, &name, &rollno)
    //     fmt.Println(strconv.Itoa(id) + ": " + name + " " + rollno)
    // }
}