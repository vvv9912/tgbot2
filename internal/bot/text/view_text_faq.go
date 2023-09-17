package text

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/botkit"
)

func ViewTextFaq() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Каталог", "https://"),
				//tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", string(sss)),
			),
		)

		msg := tgbotapi.NewMessage(botInfo.TgId, "лял")
		msg.ReplyMarkup = numericKeyboardInline
		bot.Send(msg) //todo

		return nil
	}
}
