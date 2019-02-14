package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var getAllUsers *sql.Stmt
var getUser *sql.Stmt
var getDiscordUser *sql.Stmt
var insertUser *sql.Stmt
var updateXbl *sql.Stmt
var updateName *sql.Stmt
var updateSeen *sql.Stmt
var updateTz *sql.Stmt

func mustPrepare(q string) *sql.Stmt {
	s, e := db.Prepare(q)
	if e != nil {
		panic(e)
	}
	return s
}

func mindSQL() {
	conn, err := sql.Open("mysql", sqlURI)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	db = conn
	getAllUsers = mustPrepare("SELECT `ID`,`xbl`,`destiny`,`seen`,`name`,`tz`,`discord` FROM `members` WHERE discord IS NOT NULL")
	getUser = mustPrepare("SELECT `ID`,`xbl`,`destiny`,`seen`,`name`,`tz`,`discord` FROM `members` WHERE `id`=?")
	getDiscordUser = mustPrepare("SELECT `ID`,`xbl`,`destiny`,`seen`,`name`,`tz`,`discord` FROM `members` WHERE `discord`=?")
	insertUser = mustPrepare("INSERT IGNORE INTO `members` (`discord`,`name`,`seen`,`updated_at`,`created_at`) VALUES(?,?,UNIX_TIMESTAMP(), NOW(), NOW())")
	updateXbl = mustPrepare("UPDATE `members` SET `xbl`=?, `updated_at`=NOW() WHERE `id` = ? LIMIT 1")
	updateName = mustPrepare("UPDATE `members` SET `name`=?, `updated_at`=NOW() WHERE `id` = ? LIMIT 1")
	updateSeen = mustPrepare("UPDATE `members` SET `seen`=?, `updated_at`=NOW() WHERE `id` = ? LIMIT 1")
	updateTz = mustPrepare("UPDATE `members` SET `tz`=?, `updated_at`=NOW() WHERE `id` = ? LIMIT 1")
}
