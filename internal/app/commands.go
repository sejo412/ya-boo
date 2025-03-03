package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/sejo412/ya-boo/pkg/models"
)

type Command int

const (
	CmdInit Command = iota
	CmdApprove
	CmdList
	CmdBan
)

func (c Command) String() string {
	return [...]string{"/init", "/approve", "/list", "/ban"}[c]
}

func cmdInitFirstAdmin(ctx context.Context, storage Storage, user models.User) error {
	return storage.UpsertUser(ctx, user)
}

func cmdListUsers(ctx context.Context, storage Storage) (string, error) {
	result := "|----|----------|-----------|----------|-------|\n"
	result += "| ID | Username | FirstName | LastName | Group |\n"
	result += "|----|----------|-----------|----------|-------|\n"
	users, err := storage.ListUsers(ctx)
	if err != nil {
		return "", fmt.Errorf("error list users: %w", err)
	}
	for _, user := range users {
		result += fmt.Sprintf("| %d | %s | %s | %s | %s |\n",
			user.ID, user.Username, user.FirstName, user.LastName, user.Role)
	}
	return result, nil
}

func cmdApproveUser(ctx context.Context, storage Storage, user models.User) (string, error) {
	waitingApprove, err := storage.IsWaitingApprove(ctx, user.ID)
	if err != nil {
		return "", err
	}
	if !waitingApprove {
		return "", errors.New("user is not waiting approve")
	}
	user.Role = models.RoleRegularUser
	if err = storage.UpdateUserRole(ctx, user, models.RoleRegularUser); err != nil {
		return "", fmt.Errorf("error approve user: %w", err)
	}
	return fmt.Sprintf("user successfully approved with role %s", user.Role.String()), nil
}

func cmdBanUser(ctx context.Context, storage Storage, user models.User) (string, error) {
	registered, err := storage.IsRegisteredUser(ctx, user.ID)
	if err != nil {
		return "", fmt.Errorf("error check user: %w", err)
	}
	if !registered {
		return "", errors.New("user is not registered")
	}
	err = storage.UpdateUserRole(ctx, user, models.RoleUnknown)
	if err != nil {
		return "", fmt.Errorf("error user ban: %w", err)
	}
	return fmt.Sprintf("user %d successfully banned", user.ID), nil
}
