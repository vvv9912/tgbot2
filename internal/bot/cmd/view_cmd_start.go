package cmd

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/botkit"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Hello@")); err != nil {
			return err
		}
		return nil
	}
}
