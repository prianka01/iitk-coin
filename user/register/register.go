package register

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"iitk-coin/model"

	"golang.org/x/crypto/bcrypt"
)
func addUser(user model.User) {
	 database, _ := sql.Open("sqlite3", "../../userdatabase.db")
	statement, err:= database.Prepare("INSERT INTO User (Name, Rollno, Password, Token, Access, Coins ) VALUES (?, ?, ?, ?, ?, ?)")
	if err!=nil {
		panic(err)
	}
    statement.Exec(user.Name,user.Rollno,user.Password,user.Token,user.Access,user.Coins)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; text/html; charset=utf-8")
	// w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	database, err := sql.Open("sqlite3", "../../userdatabase.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS User (Name TEXT, Rollno INTEGER PRIMARY KEY, Password TEXT, Token TEXT, Access STRING, Coins INTEGER)")
    statement.Exec()
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	rows, err := database.Query("SELECT * FROM User WHERE Rollno IN (?)",user.Rollno);
    if err!=nil {
		panic(err)
	}
	present:=false
    for rows.Next() {
		present=true
    }
	if !present {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

		if err != nil {
			res.Error = "Error While Hashing Password, Try Again"
			json.NewEncoder(w).Encode(res)
			return
		}
		user.Password = string(hash)
		user.Token=""
		user.Coins=0
		addUser(user)
		res.Result = "Registration Successful"
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Result = "Rollno already Exists!!"
	json.NewEncoder(w).Encode(res)
	return
}