package events

import "time"

type UserCreated struct {
	Nome   string
	Values any
}

func NewUserCreated() *UserCreated {
	return &UserCreated{
		Nome: "UserCreated",
	}
}

func (e *UserCreated) GetName() string {
	return e.Nome
}

func (e *UserCreated) GetPayload() any {
	return e.Values
}

func (e *UserCreated) SetPayload(values any) {
	e.Values = values
}

func (e *UserCreated) GetDateTime() time.Time {
	return time.Now()
}
