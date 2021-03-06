package records

import (
	"context"
	"database/sql"
	"encoding/json"
	"iitk-coin/model"
	"iitk-coin/pages/getdatabase"
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
	database, _ := getdatabase.GetDatabase()
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
        rows.Scan(&transaction.Type,&transaction.Sender,&transaction.Reciever,&transaction.Amount,&transaction.Tax,&transaction.TimeStamp,&transaction.Info)
        record=record+transaction.Type+" "+strconv.Itoa(transaction.Reciever)+ " from "+strconv.Itoa(transaction.Sender)+": Amount="+strconv.Itoa(transaction.Amount)+" with tax="+strconv.FormatFloat(transaction.Tax,'f',-1,64)+" at time-"+transaction.TimeStamp+" extra info="+transaction.Info;
		record+="\n";
    }
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(record)
	return
}

func FullRecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	var res model.ResponseResult
	access,err:=secretpage.HasAccess(tokenString,"gensec")
	if err!=nil {
		res.Error = "Please try again!"
		json.NewEncoder(w).Encode(res)
		return
	}
	if !access {
		res.Error = "Access to the page denied"
		json.NewEncoder(w).Encode(res)
		return
	}
	database, _ := getdatabase.GetDatabase()
	var ctx context.Context
	ctx=r.Context()
	tx, err := database.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	 if err != nil {
        res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
    }
	rows, err := database.Query("SELECT * FROM Transactions");
    if err!=nil {
		tx.Rollback()
		panic(err)
	}
    var transaction model.Transaction
	var record string
    for rows.Next() {
        rows.Scan(&transaction.Type,&transaction.Sender,&transaction.Reciever,&transaction.Amount,&transaction.Tax,&transaction.TimeStamp,&transaction.Info)
        record=record+transaction.Type+" "+strconv.Itoa(transaction.Reciever)+ " from "+strconv.Itoa(transaction.Sender)+": Amount="+strconv.Itoa(transaction.Amount)+" with tax="+strconv.FormatFloat(transaction.Tax,'f',-1,64)+" at time-"+transaction.TimeStamp+" extra info="+transaction.Info;
		record+="\n";
    }
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(record)
	return
}