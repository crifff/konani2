package main

import (
	"net/http"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/net/context"
	"konani2/service/lib/entity"
	"google.golang.org/appengine/datastore"
	"time"
	"net/url"
	"google.golang.org/appengine/taskqueue"
	"strconv"
	"google.golang.org/appengine/log"
	"fmt"
	"encoding/xml"
)

func init() {
	http.HandleFunc("/sync_calender", SyncCalender)
	http.HandleFunc("/update_title", UpdateTitle)
}

var updateTitleQueue = "default"

func UpdateTitle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	tidStr := r.FormValue("tid")
	id, err := strconv.Atoi(tidStr)
	if err != nil {
		log.Errorf(ctx, "input value: %s", tidStr)
		return
	}

	TID := entity.TID(id)

	log.Infof(ctx, "start update title TID:%d", int(TID))

	if _, err := entity.GetTitle(ctx, TID); err == datastore.ErrNoSuchEntity {
		syoboiTitle := fetchTitle(ctx, TID)
		title := &entity.Title{
			ID:            entity.TID(syoboiTitle.TitleItems.TitleItem.TID),
			Title:         syoboiTitle.TitleItems.TitleItem.Title,
			ShortTitle:    syoboiTitle.TitleItems.TitleItem.ShortTitle,
			Kana:          syoboiTitle.TitleItems.TitleItem.TitleYomi,
			Comment:       syoboiTitle.TitleItems.TitleItem.Comment,
			Category:      syoboiTitle.TitleItems.TitleItem.Cat,
			Flag:          syoboiTitle.TitleItems.TitleItem.TitleFlag,
			FirstYear:     syoboiTitle.TitleItems.TitleItem.FirstYear,
			FirstMonth:    syoboiTitle.TitleItems.TitleItem.FirstMonth,
			FirstEndYear:  syoboiTitle.TitleItems.TitleItem.FirstEndYear,
			FirstEndMonth: syoboiTitle.TitleItems.TitleItem.FirstEndMonth,
			FirstCh:       syoboiTitle.TitleItems.TitleItem.FirstCh,
			Keywords:      syoboiTitle.TitleItems.TitleItem.Keywords,
		}

		if _, err := datastore.Put(ctx, title.Key(ctx), title); err != nil {
			panic(err)
		}
	}

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

	TIDs := createTIDList(data.Items)

	tasks := createUpdateTitleTaskList(TIDs)
	var tmp []*taskqueue.Task
	for _, t := range tasks {
		tmp = append(tmp, t)
		if len(tmp) >= 100 {
			taskqueue.AddMulti(ctx, tmp, updateTitleQueue)
			tmp = make([]*taskqueue.Task, 0)

		}
	}
	if len(tmp) > 0 {
		taskqueue.AddMulti(ctx, tmp, updateTitleQueue)
		tmp = make([]*taskqueue.Task, 0)
	}
}

func createUpdateTitleTaskList(tids []entity.TID) []*taskqueue.Task {
	var tasks []*taskqueue.Task
	for _, tid := range tids {
		t := taskqueue.NewPOSTTask("/update_title", url.Values{
			"tid": {strconv.Itoa(int(tid))},
		})
		tasks = append(tasks, t)
	}
	return tasks
}

func createTIDList(items []SyoboiItem) []entity.TID {
	tmp := make(map[int]struct{})
	for _, i := range items {
		tmp[i.TID] = struct{}{}
	}

	var result []entity.TID
	for tid := range tmp {
		result = append(result, entity.TID(tid))
	}
	return result
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

func fetchTitle(ctx context.Context, tid entity.TID) TitleLookupResponse {
	client := urlfetch.Client(ctx)
	r, err := http.NewRequest("GET", fmt.Sprintf("http://cal.syoboi.jp/db.php?Command=TitleLookup&TID=%d", int(tid)), nil)
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

	var data TitleLookupResponse
	if err := xml.Unmarshal(b, &data); err != nil {
		log.Errorf(ctx, "response: %s", string(b))
		panic(err)
	}
	return data
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

type TitleLookupResponse struct {
	Result struct {
		Code    int
		Message string
	}
	TitleItems struct {
		TitleItem struct {
			TID           int
			LastUpdate    string
			Title         string
			ShortTitle    string
			TitleYomi     string
			TitleEN       string
			Comment       string
			Cat           int `xml:",attr"`
			TitleFlag     int `xml:",attr"`
			FirstYear     int `xml:",attr"`
			FirstMonth    int `xml:",attr"`
			FirstEndYear  int `xml:",attr"`
			FirstEndMonth int `xml:",attr"`
			FirstCh       string
			Keywords      string
		}
	}
}
