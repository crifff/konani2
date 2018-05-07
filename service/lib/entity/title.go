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
	Kana          string
	Comment       string `datastore:",noindex"`
	Category      int
	Flag          int
	FirstYear     int
	FirstMonth    int
	FirstEndYear  int
	FirstEndMonth int
	FirstCh       string `datastore:",noindex"`
	Keywords      string `datastore:",noindex"`
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
