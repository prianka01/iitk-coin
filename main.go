package main

import (
	"database/sql"
	"log"
	"net/http"

	"iitk-coin/user/register"

	"iitk-coin/pages/coin"
	"iitk-coin/pages/secretpage"
	"iitk-coin/user/login"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)
func main() {
	database, _ := sql.Open("sqlite3", "./userInfos.db")
 	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS User (Name TEXT, Rollno INTEGER PRIMARY KEY, Password TEXT, Token TEXT, Access BOOLEAN, Coins INTEGER)")
    statement.Exec()
	 if err!=nil {
		panic(err)
	}
    r := mux.NewRouter()
	r.HandleFunc("/signup", register.RegisterHandler).
		Methods("POST","OPTIONS")
	r.HandleFunc("/login", login.LoginHandler).
		Methods("POST","OPTIONS")
	r.HandleFunc("/secretpage", secretpage.AccessPage).
		Methods("POST","OPTIONS")
	r.HandleFunc("/awardcoins", coin.AwardCoins).
		Methods("POST","OPTIONS")
	r.HandleFunc("/getcoins", coin.GetCoins).
		Methods("GET","OPTIONS")
	r.HandleFunc("/transfercoins", coin.TransferCoins).
		Methods("POST","OPTIONS")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
    
}