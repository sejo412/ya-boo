package app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sejo412/ya-boo/pkg/models"
)

type Command int

const (
	CmdInit Command = iota
	CmdApprove
	CmdList
	CmdBan
	CmdLlmList
	CmdLlmAdd
	CmdLlmUse
	CmdLlmRemove
)

func (c Command) String() string {
	return [...]string{"/init", "/approve", "/list", "/ban", "/llmlist", "/llmadd", "/llmuse", "/llmremove"}[c]
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

func cmdLlmList(ctx context.Context, storage Storage) (string, error) {
	llmList, err := storage.GetLLMs(ctx)
	if err != nil {
		return "", fmt.Errorf("error get llmlist: %w", err)
	}
	result := "|----|------|-------------|\n"
	result += "| ID | Name | Description |\n"
	result += "|----|------|-------------|\n"
	for _, llm := range llmList {
		result += fmt.Sprintf("| %d | %s | %s |\n", llm.ID, llm.Name, llm.Description)
	}
	return result, nil
}

func cmdLlmAdd(ctx context.Context, storage Storage, message string) error {
	llm, err := parseLLM(message)
	if err != nil {
		return err
	}
	return storage.AddLLM(ctx, llm)
}

func cmdLlmUse(ctx context.Context, storage Storage, userId int64, llmId int64) error {
	llm := models.LLM{
		ID: llmId,
	}
	if err := storage.SetUserLLM(ctx, userId, llm); err != nil {
		return fmt.Errorf("error set llm for user: %w", err)
	}
	return nil
}

func cmdLlmRemove(ctx context.Context, storage Storage, id int64) error {
	llm := models.LLM{ID: id}
	if err := storage.RemoveLLM(ctx, llm); err != nil {
		return fmt.Errorf("error remove llm: %w", err)
	}
	return nil
}

func parseLLM(message string) (llm models.LLM, err error) {
	splited := strings.Split(message, " ")[1:]
	for _, v := range splited {
		vParam := strings.Split(v, "=")[0]
		vValue := strings.Split(v, "=")[1]
		switch vParam {
		case "name":
			llm.Name = vValue
		case "endpoint":
			llm.Endpoint = vValue
		case "token":
			llm.Token = vValue
		case "description":
			llm.Description = vValue
		}
	}
	if llm.Name == "" {
		return models.LLM{}, errors.New("name is required")
	}
	return llm, nil
}
