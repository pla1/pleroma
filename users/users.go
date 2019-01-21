// Dumps the users table.
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"regexp"
	"time"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=require host=%s",
		os.Getenv("PLEROMA_DB_USER"), os.Getenv("PLEROMA_DB_PASSWORD"), os.Getenv("PLEROMA_DB_NAME"), os.Getenv("PLEROMA_DB_HOST"))
	fmt.Printf("%+v\n", dbinfo)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	checkErr(err)
	fmt.Println("# Querying")
	rows, err := db.Query("SELECT id, coalesce(name,''), coalesce(nickname,''), inserted_at FROM users order by inserted_at")
	checkErr(err)

	for rows.Next() {
		var id int
		var name string
		var nickname string
		var inserted_at time.Time
		err = rows.Scan(&id, &name, &nickname, &inserted_at)
		checkErr(err)
		fmt.Printf("%10v | %50v | %50v | %25v\n", id, reg.ReplaceAllString(name, ""), reg.ReplaceAllString(nickname, ""), inserted_at)
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
