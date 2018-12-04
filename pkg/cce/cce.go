package cce

import "time"

type Limitation struct {
	CceId       string    `json:"cce_id"`		// unique cce event
	Description string    `json:"description"`	// short limitation description
	Group       string    `json:"group"`		// group submitting limitation
	Submitted   time.Time `json:"submitted"`	// time submission received
}

