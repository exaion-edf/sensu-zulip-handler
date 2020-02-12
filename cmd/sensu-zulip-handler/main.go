package main

import (
	"./zulip"
	"fmt"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu-community/sensu-plugin-sdk"
)

type HandlerConfig struct {
	sensu.PluginConfig
	ZulipUrl string
	ZulipChannel    string
	ZulipBotEmail   string
	ZulipBotKey    string
}

const (
	zulipUrl = "zulip-url"
	channel    = "channel"
	botEmail   = "bot-email"
	botKey    = "bot-key"
)

var (
	config = HandlerConfig{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-zulip-handler",
			Short:    "The Sensu Go Zulip handler for notifying a channel",
			Timeout:  10,
			Keyspace: "sensu.io/plugins/zulip/config",
		},
	}

	zulipConfigOptions = []*sensu.PluginConfigOption{
		{
			Path:      zulipUrl,
			Env:       "SENSU_ZULIP_URL",
			Argument:  zulipUrl,
			Shorthand: "u",
			Default:   "",
			Usage:     "The zulip url to send messages to, defaults to value of ZULIP_URL env variable",
			Value:     &config.ZulipUrl,
		},
		{
			Path:      channel,
			Env:       "SENSU_ZULIP_CHANNEL",
			Argument:  channel,
			Shorthand: "c",
			Default:   "#general",
			Usage:     "The channel to post messages to",
			Value:     &config.ZulipChannel,
		},
		{
			Path:      botEmail,
			Env:       "SENSU_ZULIP_BOT_EMAIL",
			Argument:  botEmail,
			Shorthand: "m",
			Default:   "",
			Usage:     "The bot that messages will be sent as",
			Value:     &config.ZulipBotEmail,
		},
		{
			Path:      botKey,
			Env:       "SENSU_ZULIP_KEY",
			Argument:  botKey,
			Shorthand: "k",
			Default:   "",
			Usage:     "The bot key",
			Value:     &config.ZulipBotKey,
		},
	}
)

func main() {
	goHandler := sensu.NewGoHandler(&config.PluginConfig, zulipConfigOptions, checkArgs, sendMessage)
	goHandler.Execute()
}

func checkArgs(_ *corev2.Event) error {
	if len(config.ZulipUrl) == 0 {
		return fmt.Errorf("--zulip-url or SENSU_ZULIP_URL environment variable is required")
	}
	if len(config.ZulipBotEmail) == 0 {
		return fmt.Errorf("--bot-email or SENSU_ZULIP_BOT_EMAIL environment variable is required")
	}
	if len(config.ZulipBotKey) == 0 {
		return fmt.Errorf("--bot-key or SENSU_ZULIP_BOT_KEY environment variable is required")
	}

	return nil
}

func messageStatus(event *corev2.Event) string {
	switch event.Check.Status {
	case 0:
		return "ok"
	case 2:
		return "alerting"
	default:
		return "warning"
	}
}

func sendMessage(event *corev2.Event) error {
	c := zulip.NewClient(
		config.ZulipUrl,
		config.ZulipBotEmail,
		config.ZulipBotKey,
	)

	msg := fmt.Sprintf("**[%s]** Check **%s** status updated for entity **%s**\n\n%s",
		messageStatus(event),
		event.Check.Name,
		event.Entity.Name,
		event.Check.Output,
		)

	_, err := c.SendMessage(
		config.ZulipChannel,
		messageStatus(event),
		msg,
	)

	return err
}
