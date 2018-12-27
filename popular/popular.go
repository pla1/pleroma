// Dumps the most favorited (liked) and boosted activity.
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
	fmt.Println("# Querying")
	rows, err := db.Query(`
		SELECT id,
       coalesce(jsonb_array_length(DATA -> 'object' -> 'likes'),0) AS quantityLikes,
			 coalesce(jsonb_array_length(data -> 'object' -> 'announcements'),0) as quantityAnnouncements,
       inserted_at
			 FROM activities
			 WHERE jsonb_typeof(DATA -> 'object' -> 'likes') = 'array'
  	 		AND DATA ->> 'type' = 'Create'
				ORDER BY id
        `)
	checkErr(err)
	var highestQuantityLiked int
	var highestIdLiked int
	var highestQuantityAnnouncements int
	var highestIdAnnouncements int
	for rows.Next() {
		var id int
		var quantityLiked int
		var quantityAnnouncements int
		var time string
		err = rows.Scan(&id, &quantityLiked, &quantityAnnouncements, &time)
		checkErr(err)
		if quantityLiked > highestQuantityLiked {
			highestQuantityLiked = quantityLiked
			highestIdLiked = id
			fmt.Println("Highest liked so far: ", highestQuantityLiked, highestIdLiked)
		}
		if quantityAnnouncements > highestQuantityAnnouncements {
			highestQuantityAnnouncements = quantityAnnouncements
			highestIdAnnouncements = id
			fmt.Println("Highest boosted quantity so far: ", highestQuantityAnnouncements, highestIdAnnouncements)
		}
	}
	printRow(db, highestIdLiked)
	printRow(db, highestIdAnnouncements)
	fmt.Printf("highestQuantityLiked quantity: %v ID: %v\n", highestQuantityLiked, highestIdLiked)
	fmt.Printf("highestQuantityBoosted quantity: %v ID: %v\n", highestQuantityAnnouncements, highestIdAnnouncements)
}

func printRow(db *sql.DB, id int) {
	stmt, err := db.Prepare("SELECT jsonb_pretty(data) as obj FROM activities where id=$1")
	checkErr(err)
	rows, err := stmt.Query(id)
	checkErr(err)
	for rows.Next() {
		var obj string
		err = rows.Scan(&obj)
		checkErr(err)
		fmt.Printf("Data: %v\n", obj)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("************ ERROR ******************")
		panic(err)
	}
}
