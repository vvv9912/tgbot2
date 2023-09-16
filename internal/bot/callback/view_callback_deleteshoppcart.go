package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/botkit"
)

func ViewCallbackdeleteshoppcart(c botkit.CorzinaStorager) botkit.ViewFunc {
	//view_callback_deleteshoppcartposition.go
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {

		err := c.DeleteCorzinaByTgID(ctx, botInfo.TgId)
		if err != nil {
			fmt.Println("DeleteCorzinaByTgID :", err)
			return err
		}
		answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Корзина удалена!"}
		bot.Send(answ)
		return nil
	}
}
