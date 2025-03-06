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
	CmdUnknown Command = iota - 1
	CmdInit
	CmdApprove
	CmdList
	CmdBan
	CmdLlmAdd
	CmdLlmRemove
	CmdLlmList
	CmdLlmUse
)

const (
	CmdInitStr      = "/init"
	CmdApproveStr   = "/approve"
	CmdListStr      = "/list"
	CmdBanStr       = "/ban"
	CmdLlmAddStr    = "/llmadd"
	CmdLlmRemoveStr = "/llmremove"
	CmdLlmListStr   = "/llmlist"
	CmdLlmUseStr    = "/llmuse"
)

func (c Command) String() string {
	return [...]string{
		CmdInitStr,
		CmdApproveStr,
		CmdListStr,
		CmdBanStr,
		CmdLlmAddStr,
		CmdLlmRemoveStr,
		CmdLlmListStr,
		CmdLlmUseStr,
	}[c]
}

func (c Command) IsAdminCommand() bool {
	return c > CmdInit && c < CmdLlmList
}

func ToCommand(s string) Command {
	switch strings.ToLower(s) {
	case CmdInitStr:
		return CmdInit
	case CmdApproveStr:
		return CmdApprove
	case CmdListStr:
		return CmdList
	case CmdBanStr:
		return CmdBan
	case CmdLlmAddStr:
		return CmdLlmAdd
	case CmdLlmRemoveStr:
		return CmdLlmRemove
	case CmdLlmListStr:
		return CmdLlmList
	case CmdLlmUseStr:
		return CmdLlmUse
	default:
		return CmdUnknown
	}
}

func cmdInitFirstAdmin(ctx context.Context, storage Storage, user models.User) error {
	return storage.UpsertUser(ctx, user)
}

func cmdListUsers(ctx context.Context, storage Storage) string {
	result := "|----|----------|-----------|----------|-------|\n"
	result += "| ID | Username | FirstName | LastName | Group |\n"
	result += "|----|----------|-----------|----------|-------|\n"
	users, err := storage.ListUsers(ctx)
	if err != nil {
		return fmt.Sprintf("error list users: %v", err)
	}
	for _, user := range users {
		result += fmt.Sprintf("| %d | %s | %s | %s | %s |\n",
			user.ID, user.Username, user.FirstName, user.LastName, user.Role)
	}
	return result
}

func cmdApproveUser(ctx context.Context, storage Storage, user models.User) string {
	waitingApprove, err := storage.IsWaitingApprove(ctx, user.ID)
	if err != nil {
		return MessageErrorApproveUser + ": " + err.Error()
	}
	if !waitingApprove {
		return MessageNotApproval
	}
	user.Role = models.RoleRegularUser
	if err = storage.UpdateUserRole(ctx, user, models.RoleRegularUser); err != nil {
		return MessageErrorApproveUser + ": " + err.Error()
	}
	return MessageSuccessApprove
}

func cmdBanUser(ctx context.Context, storage Storage, user models.User) string {
	registered, err := storage.IsRegisteredUser(ctx, user.ID)
	if err != nil {
		return MessageErrorCheckUser + ": " + err.Error()
	}
	if !registered {
		return MessageUserNotRegistered
	}
	err = storage.UpdateUserRole(ctx, user, models.RoleUnknown)
	if err != nil {
		return MessageErrorBanUser + ": " + err.Error()
	}
	return MessageSuccessBan
}

func cmdLlmList(ctx context.Context, storage Storage) string {
	llmList, err := storage.GetLLMs(ctx)
	if err != nil {
		return MessageErrorGetLLMs
	}
	result := "|----|------|-------------|\n"
	result += "| ID | Name | Description |\n"
	result += "|----|------|-------------|\n"
	for _, llm := range llmList {
		result += fmt.Sprintf("| %d | %s | %s |\n", llm.ID, llm.Name, llm.Description)
	}
	return result
}

func cmdLlmAdd(ctx context.Context, storage Storage, message string) string {
	llm, err := parseLLM(message)
	if err != nil {
		return MessageErrorParseLLM + ": " + err.Error()
	}
	if err = storage.AddLLM(ctx, llm); err != nil {
		return MessageErrorLLMAdd + ": " + err.Error()
	}
	return MessageSuccessLLMAdd
}

func cmdLlmUse(ctx context.Context, storage Storage, userID, llmID int64) string {
	llm := models.LLM{
		ID: llmID,
	}
	if err := storage.SetUserLLM(ctx, userID, llm); err != nil {
		return MessageErrorLLMUse
	}
	return MessageSuccessLLMUse
}

func cmdLlmRemove(ctx context.Context, storage Storage, id int64) string {
	llm := models.LLM{ID: id}
	if err := storage.RemoveLLM(ctx, llm); err != nil {
		return MessageErrorLLMRemove + ": " + err.Error()
	}
	return MessageSuccessLLMRemove
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
