package main

import (
	"database/sql"
	"log"
	"net/http"

	"iitk-coin/user/register"

	"iitk-coin/user/login"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)
func main() {
	database, _ := sql.Open("sqlite3", "./userData.db")
 	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS User (Name TEXT, Rollno INTEGER PRIMARY KEY, Password TEXT, Token TEXT)")
    statement.Exec()
	 if err!=nil {
		panic(err)
	}
    r := mux.NewRouter()
	r.HandleFunc("/signup", register.RegisterHandler).
		Methods("POST","OPTIONS")
	r.HandleFunc("/login", login.LoginHandler).
		Methods("POST","OPTIONS")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
    
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