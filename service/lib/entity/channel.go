package entity

import (
	"google.golang.org/appengine/datastore"
	"golang.org/x/net/context"
)

type CID int
type Channel struct {
	ID       CID
	Name     string
	URL      string
	IEPGName string
	GID      int
	Comment  string
}

func (c Channel) Key(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, "Channel", "", int64(c.ID), nil)
}
