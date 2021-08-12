package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
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

	var config struct {
		Vk struct {
			AccessToken string
		}
		Logger struct {
			Path string
		}
	}

	f, err := workfs.Open("config.toml")
	if err != nil {
		panic(err)
	}

	_, err = toml.DecodeReader(f, &config)
	if err != nil {
		panic(err)
	}
	f.Close()

	log := logger.New(logger.DefaultWriter, "MAIN")
	if config.Logger.Path != "" {
		w, err := logger.NewWriterFile(config.Logger.Path)
		if err != nil {
			panic(err)
		}
		log.SetWriter(w)
	}

	log.Info("config loaded")

	log.Info("vk authorization")
	vk, vkerr := vkapi.Auth(config.Vk.AccessToken)
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
