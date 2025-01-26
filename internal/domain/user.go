package domain

type UserID int

type User struct {
	ID        UserID
	FirstName string
	LastName  string
	Email     string
}
