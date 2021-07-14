package redeemcoins

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
	"time"
)
func createRedeemTable() error{
	database, _ := getdatabase.GetDatabase()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS RedeemRequests (Item TEXT, Sender INTEGER, Status TEXT, Info TEXT)")
    statement.Exec()
	 if err!=nil {
		panic(err)
	}
	return err
}
func addRequest(redeem model.RedeemRequest) error{
	database, _ := getdatabase.GetDatabase()
	statement, err:= database.Prepare("INSERT INTO RedeemRequests (Item, Sender, Status, Info) VALUES (?, ?, ?, ?)")
	if err!=nil {
		return err
	}
    _,err=statement.Exec(redeem.Item,redeem.Sender,redeem.Status,redeem.Info);
	return err;
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
func createTransactionTable() error{
	database, _ := getdatabase.GetDatabase()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS Transactions (Type TEXT, Sender INTEGER, Reciever INTEGER, Amount INTEGER, Tax REAL, Timestamp TEXT, Info TEXT)")
    statement.Exec()
	 if err!=nil {
		panic(err)
	}
	return err
}
func handleAcceptRequest(request model.RedeemRequest, price int) error {
	database, _ := getdatabase.GetDatabase()
	_,err:=database.Exec(`UPDATE RedeemRequests set Status=(?) AND Info=(?) WHERE Sender=(?)`,request.Status,request.Info,request.Sender)
	if err!=nil {
		return err
	}
	var transaction model.Transaction
	transaction.Type="AcceptRedeemRequest"
	transaction.Sender=request.Sender
	transaction.Reciever=1
	transaction.Amount=price
	transaction.Tax=0
	transaction.TimeStamp=time.Now().String()
	transaction.Info=request.Status+" "+request.Info
	err=createTransactionTable()
	if err!=nil {
		return err
	}
	err=addTransaction(transaction)
	return err;
}
func RedeemCoins(w http.ResponseWriter, r *http.Request) {
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
	
	var redeem model.RedeemRequest
	body, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &redeem)
	if err != nil {
		log.Fatal(err)
	}
	database, _ := getdatabase.GetDatabase()
	var ctx context.Context
	ctx=r.Context()
	tx, err := database.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	 if err != nil {
        res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
    }
	rows,err:=database.Query("SELECT Name FROM Store WHERE Name IN (?)",redeem.Item)
	if err!=nil {
		 tx.Rollback()	
		panic(err)
	}
	present:=false
    for rows.Next() {
		present=true
    }
	if !present  {
		 tx.Rollback()	
		res.Error = "No such item present"
		json.NewEncoder(w).Encode(res)
		return
	}
	redeem.Sender=rollno
	redeem.Status="pending"
	err=createRedeemTable()
	if err!=nil {
		 tx.Rollback()
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	err=addRequest(redeem)
	if err!=nil {
		 tx.Rollback()
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var transaction model.Transaction
	transaction.Type="RedeemRequest"
	transaction.Sender=user.Rollno
	transaction.Reciever=1
	transaction.Amount=0
	transaction.Tax=0
	transaction.TimeStamp=time.Now().String()
	transaction.Info=redeem.Item
	err=createTransactionTable()
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
	
	res.Result = "Request for item successful"
	json.NewEncoder(w).Encode(res)
	return
}

func AcceptRequest(w http.ResponseWriter, r *http.Request) {
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
	var request model.RedeemRequest
	body, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Fatal(err)
	}
	database, _ := getdatabase.GetDatabase()
	var ctx context.Context
	ctx=r.Context()
	tx, err := database.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	 if err != nil {
        res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
	 }
	//request would be of form itemname and rollno of sender
	//at one time only one request will be transacted
	//check for item availibilty 
	var item model.StoreItems
	itemquery,err:=database.Query("SELECT Quantity,Name,Price FROM Store WHERE Name IN (?)",request.Item)
	if err!=nil {
		 tx.Rollback()	
		panic(err)
	}
	for itemquery.Next() {
		itemquery.Scan(&item.Quantity,&item.Name, &item.Price)
    }
	if item.Quantity==0 {
		request.Info="Not enough quantity present"
		request.Status="Rejected"
		handleAcceptRequest(request,0)
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
		res.Result = "Redeem request handled"
		json.NewEncoder(w).Encode(res)
		return
	}
	//if available, check for users coins avaibility
	var sender model.User
	rows, err := database.Query("SELECT Rollno,Coins FROM User WHERE Rollno IN (?)",request.Sender);
    if err!=nil {
		panic(err)
	}
    for rows.Next() {
		rows.Scan(&sender.Rollno,&sender.Coins)
    }
	if sender.Coins-float64(item.Price)<0 {
		request.Info="Not enough coins to redeem"
		request.Status="Rejected"
		handleAcceptRequest(request,0)
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
		res.Result = "Redeem request handled"
		json.NewEncoder(w).Encode(res)
		return
	}
	//deduct coins from user in a single sql statement
	var no sql.Result
	no,err=database.Exec(`UPDATE User set Coins=Coins-(?) WHERE (Rollno=(?) AND Coins>=(?))`,float64(item.Price),request.Sender,float64(item.Price));
	x,_:=no.RowsAffected()
	if x==0 ||err != nil {
		request.Info="Error While Updating Coins or Coins not sufficient, Please try again"
		request.Status="Rejected"
		handleAcceptRequest(request,0)
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
		res.Result = "Redeem request handled"
		json.NewEncoder(w).Encode(res)
		return
	}
	//decrese item quantity
	no,err=database.Exec(`UPDATE Store set Quantity=Quantity-(?) WHERE Name=(?)`,1,request.Item);
	x,_=no.RowsAffected()
	if x==0 ||err != nil {
		tx.Rollback()
		res.Error = "Error While Updating Coins or Coins not sufficient, Please try again"
		json.NewEncoder(w).Encode(res)
		return
	}
	//if error at any point then rollback
	//change the status to redeemed
	request.Info="Succesful"
	request.Status="Approved"
	handleAcceptRequest(request,item.Price)
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	res.Result = "Redeem request handled"
	json.NewEncoder(w).Encode(res)
	return
}