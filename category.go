package chipmunk

import (
	"fmt"
	"time"
)

type category struct {
	Name   string  `json:"name"`
	Budget float64 `json:"budget"`
	Month  month   `json:"month"`
}

type month struct {
	M   time.Month `json:"m"`
	Txs []tranx    `json:"txs"`
}

func getCategory(e string) (int, error) {
	for i, c := range categories {
		if e == c.Name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("could not find category")
}

//addUser adds user to slice of users
func addCategory(c category) {
	_, err := getCategory(c.Name)
	if err != nil {
		categories = append(
			categories,
			category{
				Name:   c.Name,
				Budget: c.Budget,
			},
		)
	}
}

func sumMonth(m month) float64 {
	sum := 0.0
	for _, t := range m.Txs {
		sum += t.Cost
	}
	return sum
}
