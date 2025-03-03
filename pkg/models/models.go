package models

import "github.com/go-telegram/bot/models"

type Role int

const (
	RoleUnknown Role = iota
	RoleAdmin
	RoleRegularUser
)

type User struct {
	*models.User
	Role Role
}

func (r Role) String() string {
	return [...]string{"unknown", "admins", "regular user"}[r]
}
