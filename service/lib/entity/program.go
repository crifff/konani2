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
	Offset    int `datastore:",noindex"`
	Comment   string `datastore:",noindex"`
	Deleted   int
	Warn      int `datastore:",noindex"`
	Revision  int
	AllDay    int
}

func (p Program) Key(ctx context.Context) *datastore.Key {
	return programKey(ctx, p.ID)
}

func programKey(ctx context.Context, id PID) *datastore.Key {
	return datastore.NewKey(ctx, "Program", "", int64(id), nil)

}
