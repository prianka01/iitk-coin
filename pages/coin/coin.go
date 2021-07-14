package coin

import (
	"context"
	"database/sql"
	"encoding/json"
	"iitk-coin/model"
	"iitk-coin/pages/getdatabase"
	"iitk-coin/pages/secretpage"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	// "time"
)

type input struct {
	Rollno        		int    `json:"rollno"`
	AwardedCoins        int    `json:"awarded"`
}
type getCoin struct {
	Rollno        		int    `json:"rollno"`
	Coins        		float64    `json:"coins"`
}

type transfer struct {
	Sender      	int    `json:"sender"`
	Reciever        int    `json:"reciever"`
	Coins        	int    `json:"coins"`
}
func addTransaction(transaction model.Transaction) error{
	database, _ := getdatabase.GetDatabase()
	statement, err:= database.Prepare("INSERT INTO Transactions (Type, Sender, Reciever, Amount, Tax, Timestamp, Info) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err!=nil {
		return err
	}
    _,err=statement.Exec(transaction.Type,transaction.Sender,transaction.Reciever,transaction.Amount,transaction.Tax,transaction.TimeStamp,transaction.Info)
	return err;
}
func createTable() error{
	database, _ := getdatabase.GetDatabase()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS Transactions (Type TEXT, Sender INTEGER, Reciever INTEGER, Amount INTEGER, Tax REAL, Timestamp TEXT, Info TEXT)")
    statement.Exec()
	 if err!=nil {
		panic(err)
	}
	return err
}
func AwardCoins(w http.ResponseWriter, r *http.Request) {
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
	var user input
	body, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	database, _ := getdatabase.GetDatabase()
	var result model.User
	
	var ctx context.Context
	ctx=r.Context()
	tx, err := database.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	 if err != nil {
        res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
    }
	rows, err := database.Query("SELECT Rollno,Coins FROM User WHERE Rollno IN (?)",user.Rollno);
    if err!=nil {
		 tx.Rollback()	
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
	_,err=database.Exec(`UPDATE User set Coins= Coins+(?) WHERE Rollno =(?)`,user.AwardedCoins,user.Rollno)
	if err != nil {
		 tx.Rollback()
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var transaction model.Transaction
	transaction.Type="Awarded"
	transaction.Reciever=user.Rollno
	transaction.Amount=user.AwardedCoins
	transaction.Tax=0
	transaction.Sender=1
	transaction.TimeStamp=time.Now().String()
	transaction.Info=""
	err=createTable()
	if err!=nil {
		 tx.Rollback()
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	err=addTransaction(transaction)
	if err!=nil {
		 tx.Rollback()
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	
	res.Result = "Coins succesfully awarded"
	json.NewEncoder(w).Encode(res)
	return
}

func GetCoins(w http.ResponseWriter, r *http.Request) {
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
	var result getCoin
	var ctx context.Context
	ctx=r.Context()
	tx, err := database.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	 if err != nil {
        res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
    }
	rows, err := database.Query("SELECT Rollno,Coins FROM User WHERE Rollno IN (?)",user.Rollno);
    if err!=nil {
		tx.Rollback()
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
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(result)
	return
}

func TransferCoins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	var res model.ResponseResult
	
	rollno,err:=secretpage.StudentAccess(tokenString)
	if err!=nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	if rollno==-1 {
		res.Error = "Invalid JWT Token"
		json.NewEncoder(w).Encode(res)
		return
	}
	var request transfer
	body, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Fatal(err)
	}
	request.Sender=rollno
	database, _ := getdatabase.GetDatabase()
	var sender model.User
	var reciever model.User
	var ctx context.Context
	ctx=r.Context()
	tx, err := database.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	 if err != nil {
        res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
    }
	rows, err := database.Query("SELECT Rollno,Coins FROM User WHERE Rollno IN (?)",request.Sender);
    if err!=nil {
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
	taxtype:=2
	tax:=0.33*float64(request.Coins)
	batch_sender:=(strconv.Itoa(request.Sender))[1]-47
	batch_reciever:=(strconv.Itoa(request.Reciever))[1]-47
	if batch_sender==batch_reciever {
		tax=0.02*float64(request.Coins)
		taxtype=1
	}
	requestedcoins:=float64(request.Coins)+tax
	if sender.Coins-requestedcoins<0 {
		res.Error = "Not enough coins to transfer"
		json.NewEncoder(w).Encode(res)
		return
	}
	// time.Sleep(10*time.Second)
	var no sql.Result
	if taxtype==2 {
		no,err=database.Exec(`UPDATE User set Coins=Coins-(?) WHERE (Rollno=(?) AND Coins>=(?))`,float64(request.Coins)+0.33*float64(request.Coins),request.Sender,float64(request.Coins)+0.33*float64(request.Coins));
	}else {
		no,err=database.Exec(`UPDATE User set Coins=Coins-(?) WHERE (Rollno=(?) AND Coins>=(?))`,float64(request.Coins)+0.02*float64(request.Coins),request.Sender,float64(request.Coins)+0.02*float64(request.Coins));
	}
	x,_:=no.RowsAffected()
	if x==0 ||err != nil {
		res.Error = "Error While Updating Coins or Coins not sufficient, Please try again"
		json.NewEncoder(w).Encode(res)
		return
	}
	_,err=database.Exec(`UPDATE User set Coins=Coins+(?) WHERE Rollno=(?)`,float64(request.Coins),request.Reciever)
	if err != nil {
		tx.Rollback()
		res.Error = "Error While Updating Coins"
		json.NewEncoder(w).Encode(res)
		return
	}
	var transaction model.Transaction
	transaction.Type="Transfer"
	transaction.Sender=request.Sender
	transaction.Reciever=request.Reciever
	transaction.Amount=request.Coins
	transaction.Tax=tax
	transaction.TimeStamp=time.Now().String()
	transaction.Info=""
	err=createTable()
	if err!=nil {
		 tx.Rollback()
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	err=addTransaction(transaction)
	if err!=nil {
		 tx.Rollback()
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	res.Result = "Coins succesfully transferred"
	json.NewEncoder(w).Encode(res)
	return
}