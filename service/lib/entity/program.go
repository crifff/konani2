package entity

import (
	"time"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type PID int
type Program struct {
	ID        PID
	CID       CID
	TID       TID
	StartTime time.Time
	EndTime   time.Time
	Count     int
	Offset    int
	Comment   string
	Deleted   int
	Warn      int
	Revision  int
	AllDay    int
}

func (p Program) Key(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, "Program", "", int64(p.ID), nil)
}
