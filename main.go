package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/modules/debug"
	"github.com/Toffee-iZt/HwBot/modules/images"
	"github.com/Toffee-iZt/HwBot/modules/random"
	"github.com/Toffee-iZt/HwBot/modules/std"
	"github.com/Toffee-iZt/HwBot/modules/yalm"
	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/wfs"
)

func main() {
	println(wfs.ExecPath())
	println("PID", os.Getpid())
	println("vkapi version:", vkapi.Version)

	logWriter := logger.DefaultWriter
	if logPath := os.Getenv("LOG_PATH"); logPath != "" {
		var err error
		logWriter, err = logger.NewWriterFile(logPath)
		if err != nil {
			panic(err)
		}
	}
	log := logger.New(logWriter, "MAIN")

	log.Info("vk authorization")
	vk, vkerr := vkapi.Auth(os.Getenv("VK_TOKEN"))
	if vkerr != nil {
		log.Error("vk auth: %s", vkerr.String())
		return
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	b := bot.New(vk, logger.New(logWriter, "BOT"))
	err := b.Run(ctx, std.Setup(b), &debug.Module, &random.Module, &yalm.Module, &images.Module)
	switch err {
	case nil:
		log.Info("stopping without error")
	case context.Canceled:
		log.Info("stopping by os signal")
	default:
		log.Error("bot finished with an error: %s", err)
	}
}
