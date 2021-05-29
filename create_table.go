package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func addUser(name string, rollno string) {
	 database, _ := sql.Open("sqlite3", "./userData.db")
	statement, _ := database.Prepare("INSERT INTO user (name, rollno) VALUES (?, ?)")
    statement.Exec(name, rollno)
}
func main() {
    database, _ := sql.Open("sqlite3", "./userData.db")
 	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, name TEXT, rollno TEXT)")
    statement.Exec()
	addUser("Priyanka","190649");
	 rows, _ := database.Query("SELECT id, name, rollno FROM user")
    var id int
    var name string
    var rollno string
    for rows.Next() {
        rows.Scan(&id, &name, &rollno)
        fmt.Println(strconv.Itoa(id) + ": " + name + " " + rollno)
    }
}