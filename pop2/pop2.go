// Dumps the most favorited (liked) and boosted activity.
package main

import (
	"encoding/json"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type activityContainer struct {
	likeId int
	boostId int
	likeQuantity int
	boostQuantity int
	likeJson string
	boostJson string
}

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=require host=%s",
		os.Getenv("PLEROMA_DB_USER"), os.Getenv("PLEROMA_DB_PASSWORD"), os.Getenv("PLEROMA_DB_NAME"), os.Getenv("PLEROMA_DB_HOST"))
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()
	 now := time.Now().UTC()
	 month := int(now.Month())
	 year, week := now.ISOWeek()
	fmt.Printf("# Querying week: %v year: %v\n", week, year)
	stmt, err := db.Prepare(`
		SELECT id,
       coalesce(jsonb_array_length(DATA -> 'object' -> 'likes'),0) AS quantityLikes,
			 coalesce(jsonb_array_length(data -> 'object' -> 'announcements'),0) as quantityAnnouncements,
			 jsonb_pretty(data) as json,
       inserted_at
			 FROM activities
			 WHERE jsonb_typeof(DATA -> 'object' -> 'likes') = 'array'
  	 		AND DATA ->> 'type' = 'Create'
				and extract(year from inserted_at) = $1
				and extract(week from inserted_at) = $2
				ORDER BY id
        `)
	checkErr(err)
	rows, err := stmt.Query(year, week)
	checkErr(err)
	var container activityContainer
	for rows.Next() {
		var id int
		var quantityLiked int
		var quantityAnnouncements int
		var json string
		var time string
		err = rows.Scan(&id, &quantityLiked, &quantityAnnouncements, &json, &time)
		checkErr(err)
		if quantityLiked > container.likeQuantity {
			container.likeQuantity = quantityLiked
			container.likeId = id
			container.likeJson = json
			fmt.Println("Highest liked so far: ", container.likeQuantity, container.likeId)
		}
		if quantityAnnouncements > container.boostQuantity {
			container.boostQuantity = quantityAnnouncements
			container.boostId = id
			container.boostJson = json
			fmt.Println("Highest boosted quantity so far: ", container.boostQuantity, container.boostId)
		}
	}
	fmt.Printf("Highest Quantity Liked quantity: %v ID: %v \nJSON: %v\n", container.likeQuantity, container.likeId, container.likeJson)
	fmt.Printf("Highest Quantity Boosted quantity: %v ID: %v\nJSON: %v\n", container.boostQuantity, container.boostId, container.boostJson)
	var boostMap map[string]interface{}
  json.Unmarshal([]byte(container.boostJson), &boostMap)
	var likeMap map[string]interface{}
  json.Unmarshal([]byte(container.likeJson), &likeMap)
	fmt.Printf("Most boosted ID URL: %v\nMost liked ID URL: %v\n",boostMap["id"].(string),likeMap["id"].(string))
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("************ ERROR ******************")
		panic(err)
	}
}
