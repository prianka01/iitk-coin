package records

import (
	"context"
	"database/sql"
	"encoding/json"
	"iitk-coin/model"
	"iitk-coin/pages/secretpage"
	"log"
	"net/http"
	"strconv"
)

func ViewRecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	var res model.ResponseResult
	var user model.User
	rollno,err:=secretpage.StudentAccess(tokenString)
	if err!=nil {
		res.Error = "Please try again!"
		json.NewEncoder(w).Encode(res)
		return
	}
	if rollno==-1 {
		res.Error = "Invalid JWT Token"
		json.NewEncoder(w).Encode(res)
		return
	}
	user.Rollno=rollno
	database, err := sql.Open("sqlite3", "../../database.db")
	if err != nil {
		log.Fatal(err)
	}

	var ctx context.Context
	ctx=r.Context()
	tx, err := database.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	 if err != nil {
        res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
    }
	var result model.User
	rows, err := database.Query("SELECT Rollno FROM User WHERE Rollno IN (?)",user.Rollno);
    if err!=nil {
		tx.Rollback()
		panic(err)
	}
	present:=false
    for rows.Next() {
		present=true
		rows.Scan(&result.Rollno)
    }
	if !present {
		res.Error = "Invalid rollno"
		json.NewEncoder(w).Encode(res)
		return
	}
	rows, err = database.Query("SELECT * FROM Transactions WHERE Sender=(?) OR Reciever=(?) ",user.Rollno,user.Rollno);
    if err!=nil {
		tx.Rollback()
		panic(err)
	}
    var transaction model.Transaction
	var record string
    for rows.Next() {
        rows.Scan(&transaction.Type,&transaction.Sender,&transaction.Reciever,&transaction.Amount,&transaction.Tax)
        record=record+transaction.Type+" "+strconv.Itoa(transaction.Reciever)+ " from "+strconv.Itoa(transaction.Sender)+": Amount="+strconv.Itoa(transaction.Amount)+" with tax="+strconv.FormatFloat(transaction.Tax,'f',-1,64);
		record+="\n";
    }
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(record)
	return
}