package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var getAllUsers *sql.Stmt
var getDiscordUser *sql.Stmt
var insertUser *sql.Stmt
var insertUserMeta *sql.Stmt
var updateName *sql.Stmt
var updateSeen *sql.Stmt

func mustPrepare(q string) *sql.Stmt {
	s, e := db.Prepare(q)
	if e != nil {
		panic(e)
	}
	return s
}

func mindSQL() {
	wantMemberColumns := "m.`ID`,m.`xbl`,m.`destiny`,m.`seen`,m.`name`,m.`tz`,m.`discord`,IFNULL(mm1.meta_value, '') as `image`, IFNULL(mm2.meta_value, '') as `thumb`"
	wantMemberJoin := "members m LEFT JOIN membermeta mm1 ON(m.id=mm1.id AND mm1.meta_key='image') LEFT JOIN membermeta mm2 ON( m.id=mm2.id AND mm2.meta_key='thumb')"
	conn, err := sql.Open("mysql", sqlURI)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	db = conn
	getAllUsers = mustPrepare(fmt.Sprintf("SELECT %s FROM %s WHERE discord IS NOT NULL", wantMemberColumns, wantMemberJoin))
	getDiscordUser = mustPrepare(fmt.Sprintf("SELECT %s FROM %s WHERE `discord`=?", wantMemberColumns, wantMemberJoin))
	insertUser = mustPrepare("INSERT IGNORE INTO `members` (`discord`,`name`,`seen`,`updated_at`,`created_at`) VALUES(?,?,UNIX_TIMESTAMP(), NOW(), NOW())")
	insertUserMeta = mustPrepare("INSERT INTO `membermeta` (`member_ID`,`meta_key`,`meta_value`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `meta_value`=`meta_value`")
	updateName = mustPrepare("UPDATE `members` SET `name`=?, `updated_at`=NOW() WHERE `id` = ? LIMIT 1")
	updateSeen = mustPrepare("UPDATE `members` SET `seen`=?, `updated_at`=NOW() WHERE `id` = ? LIMIT 1")
}
