package cce

import "time"

type Limitation struct {
	Id        int       `json:"id, omitempty"`
	Limit     string    `json:"limit"`
	Group     string    `json:"group"`
	Submitted time.Time `json:"submitted, omitempty"`
}

type CCE struct {
	Title  string       `json:"cce"`
	Limits []Limitation `json:"limits"`
}
