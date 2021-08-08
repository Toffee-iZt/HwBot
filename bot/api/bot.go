package api

import (
	"HwBot/logger"
	"HwBot/vkapi"
)

//
// it is necessary to remove the need to rebuild the modules when the code changes
//

// Bot interface.
type Bot interface {
	API() *vkapi.Client
	SimpleReply(*IncomingMessage, string)

	NewCallback()
}

// Logger provides logger instance for every module.
type Logger interface {
	Info(fmt string, v ...interface{})
	Warn(fmt string, v ...interface{})
	Error(fmt string, v ...interface{})
}

// MakeLogChild ...
func MakeLogChild(log Logger, name string) Logger {
	return log.(*logger.Logger).Child(name)
}
