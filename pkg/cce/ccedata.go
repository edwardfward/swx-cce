package cce

import (
	"encoding/base64"
	"time"
)

type Record interface {
	Create() error
	Read() error
	Update() error
	Delete() error
}

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

type User struct {
	email            string
	salt             base64.Encoding
	saltedHash       base64.Encoding
	verifiedTime     time.Time
	lastActive       time.Time
	magicLink        base64.Encoding
	magicLinkExpires time.Time
}

func CreateUser(json string) *User {
	return nil
}

func (u *User) sendMagicLink() error {
	return nil
}

func (u *User) authenticateUser() (bool, error) {
	return false, nil
}

func (u *User) resetSecret() (bool, error) {
	return false, nil
}


