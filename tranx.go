package chipmunk

import "time"

type tranx struct {
	ID          int        `json:"id"`
	Cost        float64    `json:"cost"`
	Store       string     `json:"store"`
	Info        string     `json:"info"`
	Month       time.Month `json:"month"`
	User_ID     int        `json:"user_id"`
	Category_ID int        `json:"category_id"`
}
