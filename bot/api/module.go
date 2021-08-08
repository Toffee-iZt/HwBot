package api

// ModuleSymbol ...
const ModuleSymbol = "Module"

// Module struct.
type Module struct {
	Name      string
	Init      func(Bot, Logger) bool
	Terminate func()
	OnMessage func(Bot, *IncomingMessage)
	Commands  []*Command
}

// Command respresents conversation command.
type Command struct {
	Run  func(Bot, *IncomingMessage, []string)
	Cmd  string
	Desc string
	Help string
	Priv bool
	Chat bool
}
