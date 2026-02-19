package events

import "time"

// UserEmailVerified is fired when the user clicks the email link
// and the backend validates the JWT successfully.
// the handler publishes to the email.verified fanout exchange to notify all pods.
type UserEmailVerified struct {
	Nome   string
	Values any
}

func NewUserEmailVerified() *UserEmailVerified {
	return &UserEmailVerified{
		Nome: "UserEmailVerified",
	}
}

func (e *UserEmailVerified) GetName() string {
	return e.Nome
}

func (e *UserEmailVerified) GetPayload() any {
	return e.Values
}

func (e *UserEmailVerified) SetPayload(values any) {
	e.Values = values
}

func (e *UserEmailVerified) GetDateTime() time.Time {
	return time.Now()
}
