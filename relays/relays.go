// Dumps the relays found in the users table.
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=require host=%s",
		os.Getenv("PLEROMA_DB_USER"), os.Getenv("PLEROMA_DB_PASSWORD"), os.Getenv("PLEROMA_DB_NAME"), os.Getenv("PLEROMA_DB_HOST"))
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()
	fmt.Println("# Querying")
	rows, err := db.Query("select nickname, inserted_at from users where name = 'ActivityRelay'")
	checkErr(err)

	for rows.Next() {
		var nickname string
		var inserted_at time.Time
		err = rows.Scan(&nickname, &inserted_at)
		checkErr(err)
		fmt.Printf("%50v | %25v\n", nickname, inserted_at)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
