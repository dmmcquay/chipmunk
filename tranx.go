package chipmunk

import "time"

type tranx struct {
	Cost  float32    `json:"cost"`
	Store string     `json:"store"`
	Info  string     `json:"Info"`
	Month time.Month `json:"Month"`
}
