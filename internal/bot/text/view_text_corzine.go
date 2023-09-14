package text

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/botkit"
)

func ViewTextCorzine(c botkit.CorzinaStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		corz, err := c.CorzinaByTgIdwithCalalog(ctx, botInfo.TgId)
		if err != nil {
			return err
		}
		if len(corz) == 0 {
			msg := tgbotapi.NewMessage(botInfo.TgId, "Корзина пуста!")
			bot.Send(msg)
			return nil
		}
		outmsg := ""
		var numKeyInline tgbotapi.InlineKeyboardMarkup
		numKeyInline.InlineKeyboard = make([][]tgbotapi.InlineKeyboardButton, 2)

		var numKeyInline2 tgbotapi.InlineKeyboardMarkup
		numKeyInline2.InlineKeyboard = make([][]tgbotapi.InlineKeyboardButton, len(corz))
		for i := range corz {
			outmsg += fmt.Sprintf("%d товар.\nАртикул товара:%d\nНазвание:%s\n Количество:%d\n", i+1, corz[i].Article, corz[i].Name, corz[i].Quantity)
			for i := range numKeyInline.InlineKeyboard {
				var data string
				switch i {
				case 0:
					numKeyInline.InlineKeyboard[i] = make([]tgbotapi.InlineKeyboardButton, 1) //В поле создаем еще поле
					numKeyInline.InlineKeyboard[i][0].Text = "Очистить корзину"
					data = "/deleteshoppcart" //надо передавать команду +id что удалить(?) //todo
					numKeyInline.InlineKeyboard[i][0].CallbackData = &data
				case 1:
					numKeyInline.InlineKeyboard[i] = make([]tgbotapi.InlineKeyboardButton, 1) //В поле создаем еще поле
					numKeyInline.InlineKeyboard[i][0].Text = "Оформить заказ"
					data = "/placeanorder" //надо передавать команду +id что удалить(?) //todo
					numKeyInline.InlineKeyboard[i][0].CallbackData = &data
				}
			}

			for i := range numKeyInline2.InlineKeyboard {
				var data string
				numKeyInline2.InlineKeyboard[i] = make([]tgbotapi.InlineKeyboardButton, 1)
				numKeyInline2.InlineKeyboard[i][0].Text = corz[i].Name                         //
				data = fmt.Sprintf("/deleteshopposition\narticle:%d\nmsg:%d", corz[i].Article) //зачем то передавал от bot.send -> сюда |msg.id
				numKeyInline2.InlineKeyboard[i][0].CallbackData = &data

			}
		}
		msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, outmsg)
		msg1.ReplyMarkup = numKeyInline
		msgsend, err := bot.Send(msg1)
		if err != nil {
			return err
		}
		_ = msgsend // зачем
		msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "Удалить из корзины позиции:")
		msg2.ReplyMarkup = numKeyInline2
		bot.Send(msg2)
		return nil
	}
}
