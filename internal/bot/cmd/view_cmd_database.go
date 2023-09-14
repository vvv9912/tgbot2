package cmd

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tgbotv2/internal/botkit"
)

func ViewCmdAddDatabase() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {

		structMsgAddDb := botkit.BotCommand{
			Cmd: "/adddatabase",
		}
		msgAddDb, err := json.Marshal(structMsgAddDb)
		if err != nil {
			log.Println("") //todo
		}
		var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Добавить бд", string(msgAddDb)),
				//tgbotapi.NewInlineKeyboardButtonData("Добавить бд", string(sss)),
			),
		)
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Изменение базы данных")
		msg.ReplyMarkup = numericKeyboardInline
		if _, err := bot.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
