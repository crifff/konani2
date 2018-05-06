package entity

type TID int
type Title struct {
	ID         TID
	Title      string
	ShortTitle string
	Kana       string
	Comment    string
	Category   int
	Flag       int
	Urls       string
}
