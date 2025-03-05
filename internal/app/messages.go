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
	MessageLLMAddUsage    string = "use */llmadd* name=<name> endpoint=<endpoint> [token=<token>] [desc=<desc>]"
	MessageLLMAddSuccess  string = "llm add success"
	MessageLLMAddError    string = "llm add fail"
	MessageLLMError       string = "Error with llm."
	MessageUnknownCommand string = "Unknown command."
	MessageBadInitSecret  string = "Bad init secret."
	MessageNotAuthorized  string = "You are not authorized to use this command."
)
