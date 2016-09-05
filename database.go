package main

import (
	"database/sql"
	"log"
	"time"
)

var connectionPool = make(chan struct{}, 10)

func init() {
	for i := 0; i < 10; i++ {
		connectionPool <- struct{}{}
	}
}

type user struct {
	MemberID      int
	SlackID       string
	Name          string
	XBL           string
	DestinyID     string
	SeenTimestamp int64
	Seen          time.Time
	TZ            string
}

func (u *user) fromRows(row *sql.Rows) error {
	err := row.Scan(
		&u.MemberID,
		&u.SlackID,
		&u.XBL,
		&u.DestinyID,
		&u.SeenTimestamp,
		&u.Name,
		&u.TZ,
	)
	u.Seen = time.Unix(u.SeenTimestamp, 0)
	return err
}

var userCache = map[string]*user{}

type database struct {
	*sql.DB
	getUser     *sql.Stmt
	getAllUsers *sql.Stmt
	insertUser  *sql.Stmt
	updateXbl   *sql.Stmt
	updateTz    *sql.Stmt
	updateName  *sql.Stmt
}

func (d *database) user(slackID string) (*user, error) {
	var rval = &user{}
	err := d.getUser.QueryRow(slackID).Scan(
		&rval.MemberID,
		&rval.SlackID,
		&rval.XBL,
		&rval.DestinyID,
		&rval.SeenTimestamp,
		&rval.Name,
		&rval.TZ,
	)
	if err != nil {
		return nil, err
	}
	rval.Seen = time.Unix(rval.SeenTimestamp, 0)
	return rval, err
}

func (d *database) execOrLog(name string, s *sql.Stmt, args ...interface{}) {
	if _, err := s.Exec(args...); err != nil {
		log.Println(append([]interface{}{"error executing statement", name, ":", err, ":: with"}, args...))
	}
}

func (d *database) maybeInsert(slackID, name, xbl, tz string) {
	<-connectionPool
	defer func() { connectionPool <- struct{}{} }()
	_, err := d.insertUser.Exec(slackID, xbl, name, tz)
	if err != nil {
		log.Println("Error maybe inserting user:", slackID, name, xbl, tz, ":", err)
		return
	}
	user, err := d.user(slackID)
	if err != nil {
		log.Println("Error fetching user:", slackID, ":", err)
		return
	}
	if user.XBL == "" {
		d.execOrLog("xbl-update", d.updateXbl, xbl, slackID)
	}
	if user.TZ != tz {
		d.execOrLog("tx-update", d.updateTz, tz, slackID)
	}
	if user.Name != name {
		d.execOrLog("name-update", d.updateName, name, slackID)
	}
}

func (d *database) mustPrepare(sql string) *sql.Stmt {
	s, e := d.Prepare(sql)
	if e != nil {
		log.Fatal(e)
	}
	return s
}

func (d *database) prepare() {
	d.getAllUsers = d.mustPrepare("SELECT `ID`,`slack`,`xbl`,`destiny`,`seen`,`name`,`tz` FROM `members`")
	d.getUser = d.mustPrepare("SELECT `ID`,`slack`,`xbl`,`destiny`,`seen`,`name`,`tz` FROM `members` WHERE `slack`=?")
	d.insertUser = d.mustPrepare("INSERT IGNORE INTO `members` (`slack`,`xbl`,`name`,`tz`) VALUES(?,?,?,?)")
	d.updateXbl = d.mustPrepare("UPDATE `members` SET `xbl`=? WHERE `slack` = ? LIMIT 1")
	d.updateName = d.mustPrepare("UPDATE `members` SET `name`=? WHERE `slack` = ? LIMIT 1")
	d.updateTz = d.mustPrepare("UPDATE `members` SET `tz`=? WHERE `slack` = ? LIMIT 1")
}
