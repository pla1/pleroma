// Search activites.
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=require host=%s",
		os.Getenv("PLEROMA_DB_USER"), os.Getenv("PLEROMA_DB_PASSWORD"), os.Getenv("PLEROMA_DB_NAME"), os.Getenv("PLEROMA_DB_HOST"))
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()
	var searchString = os.Args[1]
	fmt.Printf("# Querying activities for: %v\n", searchString)
	rows, err := db.Query(`SELECT id, data -> 'object' ->> 'url' as url 
                       FROM activities 
                       where data ->> 'type' = 'Create' 
                       and data -> 'object' ->> 'type' = 'Note' 
                       and data -> 'object' ->> 'url' is not null 
                       and data -> 'object' ->> 'content' like $1
                       order by id desc
                       fetch first 100 rows only
                       `, searchString)
	checkErr(err)
	var id string
	var url string
	for rows.Next() {
		err = rows.Scan(&id, &url)
		checkErr(err)
		fmt.Printf("ID: %v URL: %v\n", id, url)
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
