package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "goapi:football@tcp(db:3306)/zendo"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Database is not reachable:", err)
	} else {
		fmt.Println("Successfully connected to the database!")
	}
}
