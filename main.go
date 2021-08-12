package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/bot/modules/builtin"
	"github.com/Toffee-iZt/HwBot/bot/modules/debug"
	"github.com/Toffee-iZt/HwBot/bot/modules/images"
	"github.com/Toffee-iZt/HwBot/bot/modules/random"
	"github.com/Toffee-iZt/HwBot/bot/modules/yalm"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/workfs"
)

func main() {
	println(os.Args[0])
	println(workfs.GetExec(), workfs.GetExecDir())
	println("PID", os.Getpid())
	println("vkapi version:", vkapi.Version)

	vkAccessToken := os.Getenv("VK_TOKEN")
	logPath := os.Getenv("LOG_PATH")

	log := logger.New(logger.DefaultWriter, "MAIN")
	if logPath != "" {
		w, err := logger.NewWriterFile(logPath)
		if err != nil {
			panic(err)
		}
		log.SetWriter(w)
	}

	log.Info("config loaded")

	log.Info("vk authorization")
	vk, vkerr := vkapi.Auth(vkAccessToken)
	if vkerr != nil {
		log.Error("vk auth: %s", vkerr.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	b := bot.New(vk, log.Writer())
	if !b.Run(ctx, &builtin.Module, &debug.Module, &random.Module, &yalm.Module, &images.Module) {
		cancel()
		return
	}

	<-b.Done()
	switch b.Err() {
	case nil:
		cancel()
		log.Info("stopping without error")
	case context.Canceled:
		log.Info("stopping by os signal")
	default:
		cancel()
		log.Error("bot finished with an error: %s", b.Err())
	}
}
