package main

import (
	"net/http"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"konani2/service/lib/entity"
	"google.golang.org/appengine/datastore"
	"time"
)

func init() {
	http.HandleFunc("/sync_calender", SyncCalender)
}

func SyncCalender(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	data := fetchCalendar(ctx)

	//channel
	channels := createChannels(data.Channels)
	var cKeys []*datastore.Key
	for _, c := range channels {
		cKeys = append(cKeys, c.Key(ctx))
	}
	if _, err := datastore.PutMulti(ctx, cKeys, channels); err != nil {
		panic(err)
	}
	//channel
	programs := createPrograms(data.Items)
	var pKeys []*datastore.Key
	for _, c := range programs {
		pKeys = append(pKeys, c.Key(ctx))
	}
	if _, err := datastore.PutMulti(ctx, pKeys, programs); err != nil {
		panic(err)
	}

	fmt.Printf("%#v", data)
}

func createPrograms(items []SyoboiItem) []*entity.Program {
	var result []*entity.Program
	for _, i := range items {
		result = append(result, &entity.Program{
			ID:        entity.PID(i.PID),
			CID:       entity.CID(i.ChID),
			TID:       entity.TID(i.TID),
			StartTime: time.Unix(int64(i.StTime), 0),
			EndTime:   time.Unix(int64(i.EdTime), 0),
			Count:     i.Count,
			Offset:    i.StOffset,
			Comment:   i.ProgComment,
			Deleted:   i.Deleted,
			Warn:      i.Warn,
			Revision:  i.Revision,
			AllDay:    i.AllDay,
		})
	}
	return result
}
func createChannels(channels map[int]SyoboiChannel) []*entity.Channel {
	var result []*entity.Channel
	for _, c := range channels {
		result = append(result, &entity.Channel{
			ID:       entity.CID(c.ChID),
			Name:     c.ChName,
			URL:      c.ChURL,
			IEPGName: c.ChiEPGName,
			GID:      c.ChGID,
			Comment:  c.ChComment,
		})
	}
	return result
}

func fetchCalendar(ctx context.Context) SyoboiResponse {
	client := urlfetch.Client(ctx)
	r, err := http.NewRequest("GET", "http://cal.syoboi.jp/rss2.php?alt=json&days=7", nil)
	if err != nil {
		panic(err)
	}
	response, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var data SyoboiResponse
	if err := json.Unmarshal(b, &data); err != nil {
		panic(err)
	}
	return data
}

type SyoboiResponse struct {
	Items    []SyoboiItem          `json:"items"`
	Channels map[int]SyoboiChannel `json:"chInfo"`
}

type SyoboiItem struct {
	StTime      int `json:",string"`
	EdTime      int `json:",string"`
	LastUpdate  string
	Count       int `json:",string"`
	StOffset    int `json:",string"`
	TID         int `json:",string"`
	PID         int `json:",string"`
	ProgComment string
	ChID        int `json:",string"`
	SubTitle    string
	Flag        int `json:",string"`
	Deleted     int `json:",string"`
	Warn        int `json:",string"`
	Revision    int `json:",string"`
	AllDay      int `json:",string"`
	Title       string
	ShortTitle  string
	Cat         int `json:",string"`
	Urls        string
	ChName      string
	ChURL       string
	ChGID       int `json:",string"`
}

type SyoboiChannel struct {
	ChID       int `json:",string"`
	ChName     string
	ChURL      string
	ChiEPGName string
	ChGID      int `json:",string"`
	ChComment  string
}
