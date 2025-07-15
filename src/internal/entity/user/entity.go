package user

import "time"

type UserToRegiser struct {
	ID        int64
	FirstName string
	LastName  string
	BirthDate time.Time
	RegisterStep
}
