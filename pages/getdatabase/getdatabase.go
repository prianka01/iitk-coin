package getdatabase

import (
	"database/sql"
	"log"
)

func GetDatabase() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "../../datarecord.db")
	if err != nil {
		log.Fatal(err)
	}
	return database, err
}