package models

import "github.com/go-telegram/bot/models"

type LLM struct {
	ID          int64
	Name        string
	Endpoint    string
	Token       string
	Description string
}

type Role int

const (
	RoleUnknown Role = iota
	RoleAdmin
	RoleRegularUser
)

type User struct {
	*models.User
	Role Role
	LLM  LLM
}

func (r Role) String() string {
	return [...]string{"unknown", "admins", "regular user"}[r]
}
