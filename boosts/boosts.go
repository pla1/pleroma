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
        rows, err := db.Query("SELECT id, coalesce(jsonb_array_length(data -> 'object' -> 'announcements'),0) as quantity FROM activities where data ->> 'type' = 'Create'")
        checkErr(err)
	var highest int
        var high_id int
        for rows.Next() {
            var id int
            var quantity int
            err = rows.Scan(&id, &quantity)
            if quantity > highest {
              highest = quantity
              high_id = id
              fmt.Println("High so far: ", highest, high_id)
	    }
            checkErr(err)
        }
	fmt.Printf("Highest quantity: %v ID: %v\n", highest, high_id)
    }

    func checkErr(err error) {
        if err != nil {
            panic(err)
        }
    }
