package main

import (
	"HwBot/bot"
	"HwBot/common/execdir"
	"HwBot/logger"
	"HwBot/vkapi"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
)

func main() {
	println(os.Args[0])
	println(execdir.GetExec(), execdir.GetExecDir())
	println("PID", os.Getpid())
	println("vkapi version:", vkapi.Version)

	var config struct {
		Bot bot.Config
		Vk  struct {
			AccessToken string
		}
		Logger struct {
			Path string
		}
	}

	f, err := execdir.Open("config.toml")
	if err != nil {
		panic(err)
	}

	_, err = toml.DecodeReader(f, &config)
	if err != nil {
		panic(err)
	}
	f.Close()

	var logWriter *logger.Writer

	if config.Logger.Path != "" {
		logWriter, err = logger.NewWriterFile(config.Logger.Path)
		if err != nil {
			panic(err)
		}
	} else {
		logWriter = logger.NewWriter(os.Stderr, true)
	}

	log := logger.New(logWriter, "MAIN")
	log.SetWriter(logWriter)

	log.Info("config loaded")

	log.Info("vk authorization")
	vk, vkerr := vkapi.Auth(config.Vk.AccessToken)
	if vkerr != nil {
		log.Error("vk auth: %s", vkerr.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	b := bot.New(vk, config.Bot, logWriter)
	if !b.Run(ctx) {
		cancel()
		return
	}

	<-b.Done()
	switch b.Err() {
	case nil:
		log.Info("nil error stop")
	case context.Canceled:
		log.Info("stopping by os signal")
	default:
		cancel()
		log.Error("bot finished with an error: %s", b.Err())
	}
}
