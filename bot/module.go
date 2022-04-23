package bot

import (
	"github.com/Toffee-iZt/HwBot/common"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
)

// Module struct.
type Module struct {
	Name      string
	Init      func() error
	Terminate func()
	Callback  func(ctx *Callback, payload vkapi.JSONData)
	Commands  []*Command
	log       *logger.Logger
}

// Prefixes are the characters with which commands must begin.
const Prefixes = "/!"

// Command respresents conversation command.
type Command struct {
	Run         func(*Context, *Message)
	Cmd         []string
	Description string
	Help        string
	Options     common.Flag
}

// Command options.
const (
	OptionInDialog common.Flag = 1 << iota
	OptionInChat
)

// Message type.
type Message struct {
	vkapi.Message
	Args []string
}
