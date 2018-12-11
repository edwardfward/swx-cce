package access

import (
	"time"
)

type UserLink struct {
	linkUrl string
	issued  time.Time
	expires time.Time
}

type User struct {
	email             string    // valid email address
	role              string    // user's role
	hashedSecret      string    // bcrypt generated salt and hash of password
	resetLink		  UserLink	// reset user's password
	verifyLink		  UserLink  // email verification link
	magicLink		  UserLink  // magic link
}



