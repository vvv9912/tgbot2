package text

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/botkit"
	"tgbotv2/internal/model"
)

func ViewTextOrder(o botkit.OrderStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		order, err := o.OrdersByTgID(ctx, botInfo.TgId)
		if err != nil {
			fmt.Println("err ", err)
			return err
		}
		if len(order) == 0 {
			msg := tgbotapi.NewMessage(botInfo.TgId, "Заказов нет!")
			bot.Send(msg)
			return nil
		}
		outmsg := ""
		for i := range order {
			var Corz []model.OrderCorz
			err := json.Unmarshal([]byte(order[i].Order), &Corz)
			if err != nil {
				fmt.Println("decode order json, err: ", err)
				return err
			}
			//Проверка на корз = 0
			outmsg += fmt.Sprintf("➖➖➖➖➖➖➖\nID Заказа №%d\n", order[i].ID)
			for k := range Corz {
				outmsg += fmt.Sprintf("%d товар.\nАртикул товара: %d\nНазвание: %s\nКоличество: %d\nЦена: %0.2f\n", k+1, Corz[k].Article, Corz[k].Name, Corz[k].Quantity, Corz[k].Price)
			}
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, outmsg)
		bot.Send(msg)

		return nil
	}
}
