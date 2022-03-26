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
	println(wfs.GetExecName(), wfs.GetExecDir())
	println("PID", os.Getpid())
	println("vkapi version:", vkapi.Version)

	log := logger.New(logger.DefaultWriter, "MAIN")
	if logPath := os.Getenv("LOG_PATH"); logPath != "" {
		w, err := logger.NewWriterFile(logPath)
		if err != nil {
			panic(err)
		}
		log.SetWriter(w)
	}

	log.Info("vk authorization")
	vk, vkerr := vkapi.Auth(os.Getenv("VK_TOKEN"))
	if vkerr != nil {
		log.Error("vk auth: %s", vkerr.Error())
		return
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	b := bot.New(vk, log.Copy("BOT"))
	err := b.Run(ctx, true, std.Setup(b), &debug.Module, &random.Module, &yalm.Module, &images.Module)
	switch err {
	case nil:
		log.Info("stopping without error")
	case context.Canceled:
		log.Info("stopping by os signal")
	default:
		log.Error("bot finished with an error: %s", err)
	}
}
