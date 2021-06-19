package coin

import (
	"database/sql"
	"encoding/json"
	"iitk-coin/model"

	// "iitk-coin/packages/userip"
	"io/ioutil"
	"log"
	"net/http"
)

type input struct {
	Rollno        		int    `json:"rollno"`
	AwardedCoins        int    `json:"awarded"`
}
type getCoin struct {
	Rollno        		int    `json:"rollno"`
	Coins        		int    `json:"coins"`
}

type transfer struct {
	Sender      	int    `json:"sender"`
	Reciever        int    `json:"reciever"`
	Coins        		int    `json:"coins"`
}
func AwardCoins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user input
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	database, err := sql.Open("sqlite3", "../../userInfos.db")

	if err != nil {
		log.Fatal(err)
	}
	var result model.User
	var res model.ResponseResult
	// var ctx context.Context
	// tx,err:= database.BeginTx(ctx,nil);
	//  if err != nil {
    //     res.Error = err.Error()
	// 	json.NewEncoder(w).Encode(res)
    // }
	rows, err := database.Query("SELECT Rollno,Coins FROM User WHERE Rollno IN (?)",user.Rollno);
    if err!=nil {
		// tx.Rollback()	
		panic(err)
	}
	present:=false
    for rows.Next() {
		present=true
		rows.Scan(&result.Rollno,&result.Coins)
    }
	if !present {
		res.Error = "Invalid rollno"
		json.NewEncoder(w).Encode(res)
		return
	}
	_,err=database.Exec(`UPDATE User set Coins=(?) WHERE Rollno=(?)`,result.Coins+user.AwardedCoins,result.Rollno)
	if err != nil {
		// tx.Rollback()
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	// tx.Commit()
	res.Result = "Coins succesfully awarded"
	json.NewEncoder(w).Encode(res)
	return
}

func GetCoins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	database, err := sql.Open("sqlite3", "../../userInfos.db")

	if err != nil {
		log.Fatal(err)
	}
	var result getCoin
	var res model.ResponseResult
	// var ctx context.Context
	// tx,err:= database.BeginTx(ctx,nil);
	//  if err != nil {
    //     res.Error = err.Error()
	// 	json.NewEncoder(w).Encode(res)
    // }
	rows, err := database.Query("SELECT Rollno,Coins FROM User WHERE Rollno IN (?)",user.Rollno);
    if err!=nil {
		// tx.Rollback()
		panic(err)
	}
	present:=false
    for rows.Next() {
		present=true
		rows.Scan(&result.Rollno,&result.Coins)
    }
	if !present {
		res.Error = "Invalid rollno"
		json.NewEncoder(w).Encode(res)
		return
	}
	// tx.Commit()
	json.NewEncoder(w).Encode(result)
	return
}

func TransferCoins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request transfer
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Fatal(err)
	}
	database, err := sql.Open("sqlite3", "../../userInfos.db")

	if err != nil {
		log.Fatal(err)
	}
	var sender model.User
	var reciever model.User
	var res model.ResponseResult
	// var ctx context.Context
	//  userIP, err := userip.FromRequest(r)
    // if err != nil {
    //     http.Error(w, err.Error(), http.StatusBadRequest)
    //     return
    // }
    // ctx = userip.NewContext(ctx, userIP)
	// tx,err:= database.BeginTx(ctx,nil);
	//  if err != nil {
    //     res.Error = err.Error()
	// 	json.NewEncoder(w).Encode(res)
    // }
	rows, err := database.Query("SELECT Rollno,Coins FROM User WHERE Rollno IN (?)",request.Sender);
    if err!=nil {
		// tx.Rollback()
		panic(err)
	}
	present:=false
    for rows.Next() {
		present=true
		rows.Scan(&sender.Rollno,&sender.Coins)
    }
	if !present {
		res.Error = "Invalid sender rollno"
		json.NewEncoder(w).Encode(res)
		return
	}
	rows, err = database.Query("SELECT Rollno,Coins FROM User WHERE Rollno IN (?)",request.Reciever);
    if err!=nil {
		// tx.Rollback()
		panic(err)
	}
	present=false
    for rows.Next() {
		present=true
		rows.Scan(&reciever.Rollno,&reciever.Coins)
    }
	if !present {
		res.Error = "Invalid reciever rollno"
		json.NewEncoder(w).Encode(res)
		return
	}
	if sender.Coins-request.Coins<0 {
		res.Error = "Not enough coins to transfer"
		json.NewEncoder(w).Encode(res)
		return
	}
	_,err=database.Exec(`UPDATE User set Coins=(?) WHERE Rollno=(?)`,sender.Coins-request.Coins,sender.Rollno)
	if err != nil {
		// tx.Rollback()
		res.Error = "Error While Updating Coins"
		json.NewEncoder(w).Encode(res)
		return
	}
	_,err=database.Exec(`UPDATE User set Coins=(?) WHERE Rollno=(?)`,reciever.Coins+request.Coins,reciever.Rollno)
	if err != nil {
		// tx.Rollback()
		res.Error = "Error While Updating Coins"
		json.NewEncoder(w).Encode(res)
		return
	}
	// tx.Commit()
	res.Result = "Coins succesfully transferred"
	json.NewEncoder(w).Encode(res)
	return
}