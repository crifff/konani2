package entity

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type TID int
type Title struct {
	ID            TID
	Title         string
	ShortTitle    string
	Comment       string `datastore:",noindex"`
	Category      int
	Flag          int
	Urls      string `datastore:",noindex"`
}

func (t Title) Key(ctx context.Context) *datastore.Key {
	return titleKey(ctx, t.ID)
}

func titleKey(ctx context.Context, id TID) *datastore.Key {
	return datastore.NewKey(ctx, "Title", "", int64(id), nil)
}

func GetTitle(ctx context.Context, id TID) (*Title, error) {
	var result *Title
	key := titleKey(ctx, id)
	err := datastore.Get(ctx, key, result)
	return result, err
}
