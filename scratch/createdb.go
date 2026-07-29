package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	connString := "server=100.109.65.2;user id=sa;password=31052006Hoa*;port=1433;database=master"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err.Error())
	}

	fmt.Println("Connected to master DB. Creating VisualFinanceDB...")
	
	// Create Database
	_, err = db.Exec("CREATE DATABASE VisualFinanceDB")
	if err != nil {
		fmt.Println("Database might already exist or error:", err.Error())
	} else {
		fmt.Println("Database VisualFinanceDB created successfully!")
	}
}
