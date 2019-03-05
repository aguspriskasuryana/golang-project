package main

import (
    "database/sql"
)

func Conn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := ""
    dbName := "db"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}
