package entity

type TermsOfService struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
	TimeStamp
}
