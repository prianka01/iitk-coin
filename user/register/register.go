package register

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/prianka01/iitk-coin/model"

	"golang.org/x/crypto/bcrypt"
)
func addUser(user model.User) {
	 database, _ := sql.Open("sqlite3", "../../userData.db")
	statement, _ := database.Prepare("INSERT INTO user (Name, Rollno, Password, Token) VALUES (?, ?, ?, ?)")
    statement.Exec(user.Name,user.Rollno,user.Password,user.Token)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	database, err := sql.Open("sqlite3", "../../userData.db")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	_, err = database.Query("SELECT * FROM user WHERE Rollno = (?)",user.Rollno);

	if err != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

		if err != nil {
			res.Error = "Error While Hashing Password, Try Again"
			json.NewEncoder(w).Encode(res)
			return
		}
		user.Password = string(hash)
		user.Token=""
		addUser(user)
		if err != nil {
			res.Error = "Error While Creating User, Try Again"
			json.NewEncoder(w).Encode(res)
			return
		}
		res.Result = "Registration Successful"
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Result = "Username already Exists!!"
	json.NewEncoder(w).Encode(res)
	return
}