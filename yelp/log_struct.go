package yelp

import "time"

type Log struct {
	Id         int
	ZipCode    int
	Alias      string
	Limit      int
	Offset     int
	IsDone     bool
	Error      string
	CreatedAt  time.Time
	ModifiedAt time.Time
}
