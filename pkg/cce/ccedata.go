package cce

import "time"

type Limit struct {
	Id        int       `json:"id, omitempty"`
	Limit     string    `json:"limit"`
	Group     string    `json:"group"`
	Submitted time.Time `json:"submitted, omitempty"`
}

type CCE struct {
	Title  string  `json:"cce"`
	Limits []Limit `json:"limits"`
}
