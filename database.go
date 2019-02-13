package main

import (
	"database/sql"
	"time"
)

type user struct {
	MemberID      int
	Name          string
	XBL           string
	DestinyID     string
	SeenTimestamp int64
	Seen          time.Time
	TZ            string
	DiscordID     string
}

func (u *user) fromRows(row *sql.Rows) error {
	err := row.Scan(
		&u.MemberID,
		&u.XBL,
		&u.DestinyID,
		&u.SeenTimestamp,
		&u.Name,
		&u.TZ,
		&u.DiscordID,
	)
	u.Seen = time.Unix(u.SeenTimestamp, 0)
	return err
}
