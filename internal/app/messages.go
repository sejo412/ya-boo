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
	MessageUnknownCommand string = `
Unknown command.
`
	MessageBadInitSecret string = `
Bad init secret.
`
)
