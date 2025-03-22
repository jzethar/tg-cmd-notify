package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"

	telebot "gopkg.in/telebot.v4"
)

type Config struct {
	ChatId     int    `yaml:"chat_id"`
	TgBotToken string `yaml:"tg_bot_token"`
	Instance   string `yaml:"instance"`
	ThreadID   int    `yaml:"thread_id"`
}

func newConfig(configFilePath string) *Config {
	var err error
	var file *os.File
	var data []byte
	var config *Config

	if file, err = os.Open(configFilePath); err != nil {
		log.Fatal().Msg("can't open file")
	}
	defer file.Close()

	if data, err = io.ReadAll(file); err != nil {
		log.Fatal().Msg("can't read file")
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Fatal().Msg("can't unmarshal file")
	}
	return config
}

func newBot(config *Config) *telebot.Bot {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  config.TgBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal().Msgf("Failed to create bot: %v", err)
	}
	return bot
}

func notifyInTG(bot *telebot.Bot, config *Config, sentArgs []string) {
	var (
		err    error
		buffer bytes.Buffer
	)
	buffer.WriteString("On instance: ")
	buffer.WriteString(config.Instance)
	buffer.WriteString(" finished cmd: ")
	for i := 1; i < len(os.Args); i++ {
		buffer.WriteString(sentArgs[i])
		buffer.WriteString(" ")
	}

	if config.ThreadID == -1 {
		_, err = bot.Send(telebot.ChatID(config.ChatId), buffer.String())
	} else {
		_, err = bot.Send(telebot.ChatID(config.ChatId), buffer.String(), &telebot.SendOptions{
			ThreadID: config.ThreadID,
		})
	}

	if err != nil {
		log.Fatal().Msgf("Failed to send message: %v", err)
	}

}

func main() {
	var (
		err     error
		cmdArgs []string
		cmdProg string
	)
	notifyConfig := os.Getenv("NOTIFY_CONFIG")
	if len(notifyConfig) == 0 {
		log.Fatal().Msg("No config provided in NOTIFY_CONFIG")
	}
	config := newConfig(notifyConfig)
	bot := newBot(config)

	switch len(os.Args) {
	case 1:
		return
	case 2:
		cmdProg = os.Args[1]
	default:
		if len(os.Args) > 2 {
			cmdProg = os.Args[1]
			cmdArgs = os.Args[2:]
		} else {
			return
		}
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer func() {
		log.Info().Msg("Finishing service and sending to tg")
		cancel()
		notifyInTG(bot, config, os.Args)
		log.Info().Msg("Message sent")
	}()
	cmd := exec.CommandContext(ctx, cmdProg, cmdArgs...)
	cmd.Stdout = os.Stdout

	if err = cmd.Run(); err != nil {
		log.Info().Msgf("could not run command: %v", err)
	}
}
