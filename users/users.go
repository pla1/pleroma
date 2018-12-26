package main

    import (
        "database/sql"
        "fmt"
        _ "github.com/lib/pq"
        "time"
        "os"
    )

    func main() {
        dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=require host=%s",
            os.Getenv("PLEROMA_DB_USER"), os.Getenv("PLEROMA_DB_PASSWORD"), os.Getenv("PLEROMA_DB_NAME"), os.Getenv("PLEROMA_DB_HOST"))
        fmt.Printf("%+v\n", dbinfo)
        db, err := sql.Open("postgres", dbinfo)
        checkErr(err)
        defer db.Close()

        fmt.Println("# Querying")
        rows, err := db.Query("SELECT id, coalesce(name,''), coalesce(nickname,''), inserted_at FROM users")
        checkErr(err)

        for rows.Next() {
            var id int
            var name string
            var nickname string
            var inserted_at time.Time
            err = rows.Scan(&id, &name, &nickname, &inserted_at)
            checkErr(err)
            fmt.Println("id | name | nickname | inserted_at")
            fmt.Printf("%3v | %8v | %6v | %6v\n", id, name, nickname, inserted_at)
        }

    }

    func checkErr(err error) {
        if err != nil {
            panic(err)
        }
    }
