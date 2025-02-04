package app

const (
	MessageInit string = `
Bot in *INIT* mode.
Type /init <secret> to initialize the bot with first admin:
*ID*: %d
*Username*: %s
*First name*: %s
*Last name*: %s
`
	MessageInitOk string = `
Bot initialized successfully.
Switch to normal mode.
`
	MessageHelper string = `
LLM not chosen
use:
*/llmlist* - list available llm
*/llmuse*  - use llm by id
`
	MessageErrorCheckAdmin   string = "error checking admin role"
	MessageErrorCheckUser    string = "error checking user"
	MessageErrorApproveUser  string = "error approving user"
	MessageErrorBanUser      string = "error ban user"
	MessageErrorGetLLMs      string = "error getting llms"
	MessageErrorParseLLM     string = "error parse llm"
	MessageErrorLLMAdd       string = "llm add fail"
	MessageErrorLLMUse       string = "error llm use"
	MessageErrorLLMRemove    string = "error llm remove"
	MessageNotApproval       string = "user not waiting approve"
	MessageSuccessApprove    string = "user successfully approved"
	MessageSuccessBan        string = "user successfully banned"
	MessageSuccessLLMAdd     string = "llm successfully added"
	MessageSuccessLLMRemove  string = "llm successfully removed"
	MessageSuccessLLMUse     string = "llm switch success"
	MessageLLMAddUsage       string = "use */llmadd* name=<name> endpoint=<endpoint> [token=<token>] [desc=<desc>]"
	MessageLLMInternalError  string = "internal error with llm"
	MessageUnknownCommand    string = "unknown command"
	MessageBadInitSecret     string = "bad init secret"
	MessageNotAuthorized     string = "you are not authorized to use this command"
	MessageUserNotRegistered string = "user not registered"
)
