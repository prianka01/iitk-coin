package main

import (
	"database/sql"
	"log"
	"net/http"

	"iitk-coin/user/register"

	"iitk-coin/pages/coin"
	"iitk-coin/pages/listitems"
	"iitk-coin/pages/redeemcoins"
	"iitk-coin/pages/secretpage"
	"iitk-coin/user/login"
	"iitk-coin/user/records"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)
func main() {
	database, _ := sql.Open("sqlite3", "./datarecord.db")
 	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS User (Name TEXT, Rollno INTEGER PRIMARY KEY, Password TEXT, Token TEXT, Access STRING, Coins REAL, Events INTEGER)")
    statement.Exec()
	 if err!=nil {
		panic(err)
	}
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS Transactions (Type TEXT, Sender INTEGER, Reciever INTEGER, Amount INTEGER, Tax REAL, Timestamp TEXT)")
    statement.Exec()
	 if err!=nil {
		panic(err)
	}
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS Store (Indice INTEGER, Name TEXT, Quantity INTEGER, Price INTEGER)")
    statement.Exec()
	 if err!=nil {
		panic(err)
	}
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS RedeemRequests (Item TEXT, Sender INTEGER, Status TEXT, Info TEXT)")
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
	r.HandleFunc("/records", records.ViewRecords).
		Methods("POST","OPTIONS")
	r.HandleFunc("/fullrecords", records.FullRecords).
		Methods("POST","OPTIONS")
	r.HandleFunc("/additem", listitems.AddItem).
		Methods("POST","OPTIONS")
	r.HandleFunc("/redeem",redeemcoins.RedeemCoins).
		Methods("POST","OPTIONS")
	r.HandleFunc("/acceptredeem",redeemcoins.AcceptRequest).
		Methods("POST","OPTIONS")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
    
}