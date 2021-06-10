package login

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"iitk-coin/model"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	database, err := sql.Open("sqlite3", "../../userData.db")

	if err != nil {
		log.Fatal(err)
	}
	var result model.User
	var res model.ResponseResult

	rows, err := database.Query("SELECT * FROM User WHERE Rollno IN (?)",user.Rollno);
    if err!=nil {
		panic(err)
	}
	present:=false
    for rows.Next() {
		present=true
		rows.Scan(&result.Name,&result.Rollno,&result.Password,&result.Token)
    }

	if !present {
		res.Error = "Invalid rollno"
		json.NewEncoder(w).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		res.Error = "Invalid password"
		json.NewEncoder(w).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":  result.Name,
		"rollno":result.Rollno,

	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		res.Error = "Error while generating token,Try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""
	json.NewEncoder(w).Encode(result)
}