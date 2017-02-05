package chipmunk

import "time"

type tranx struct {
	ID          int       `json:"id"`
	Cost        float64   `json:"cost"`
	Store       string    `json:"store"`
	Info        string    `json:"info"`
	Date        time.Time `json:"date"`
	User_ID     int       `json:"user_id"`
	Category_ID int       `json:"category_id"`
}
