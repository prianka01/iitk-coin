package listitems

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
)

func addIInStore(item model.StoreItems) error{
	database, _ := getdatabase.GetDatabase()
	statement, err:= database.Prepare("INSERT INTO Store (Name, Quantity, Price) VALUES (?, ?, ?)")
	if err!=nil {
		return err
	}
    _,err=statement.Exec(item.Name,item.Quantity,item.Price)
	return err;
}
func createTable() error{
	database, _ := getdatabase.GetDatabase()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS Store (Indice INTEGER, Name TEXT, Quantity INTEGER, Price INTEGER)")
    _,err=statement.Exec()
	 if err!=nil {
		panic(err)
	}
	return err
}
func AddItem(w http.ResponseWriter, r *http.Request) {
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
	var request model.StoreItems
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
	err=createTable()
	if err!=nil {
		 tx.Rollback()
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	rows, err := database.Query("SELECT Indice,Name,Quantity FROM Store WHERE Name IN (?)",request.Name);
    if err!=nil {
		 tx.Rollback()	
		panic(err)
	}
	present:=false
    for rows.Next() {
		present=true
		rows.Scan(&request.Name,&request.Quantity)
    }
	if !present {
		// rows,_=database.Query("SELECT * FROM Store")
		// var co=0
		// for rows.Next() {
		// 	co=co+1
    	// }
		// request.Index=1
		err=addIInStore(request)
	}
	if present {
		_,err=database.Exec(`UPDATE Store set Quantity=Quantity+(?) WHERE Name=(?)`,request.Quantity,request.Name)
	}
	if err!=nil {
		tx.Rollback()
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	
	res.Result = "Item Added in Store"
	json.NewEncoder(w).Encode(res)
	return
}