package callback_order

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/bot/constant"
	"tgbotv2/internal/botkit"
	"tgbotv2/internal/model"
	"time"
)

func ViewCallbackAddOrderSamovivoz(o botkit.OrderStorager, c botkit.CorzinaStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		//выбрать доставку https://github.com/vvv99Корзина пуста!12/mytgbot/blob/main/tg2/parscommand.go#L714
		//написать адм или выбрать пвз и тп...
		//потом
		corz, err := c.CorzinaByTgIdwithCalalog(ctx, botInfo.TgId)
		//corz, err := c.CorzinaByTgId(ctx, botInfo.TgId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if len(corz) == 0 {
			msg := tgbotapi.NewMessage(botInfo.TgId, "Корзина пуста!")
			bot.Send(msg)
			return nil
		}
		////
		//_ = corz
		//type corzfororder model.OrderCorz
		//запрос с ценой и корзиной
		corzForOrder := make([]model.OrderCorz, len(corz))
		for i := range corz {
			corzForOrder[i].Article = corz[i].Article
			corzForOrder[i].Quantity = corz[i].Quantity
			corzForOrder[i].Price = corz[i].Price
			corzForOrder[i].ID = corz[i].ID
			corzForOrder[i].Name = corz[i].Name
		}
		msgcorz, err := json.Marshal(corzForOrder)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = o.AddOrders(ctx, model.Orders{
			//ID:          0,
			TgID:        botInfo.TgId,
			StatusOrder: 0,
			Pvz:         "{}",
			Order:       string(msgcorz),

			CreatedAt:    time.Now().UTC().Add(3 * time.Hour),
			TypeDostavka: constant.D_SAMOVIVOZ,
			//ReadAt:      ,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		msg := tgbotapi.NewMessage(botInfo.TgId, "Заказ добавлен!")
		bot.Send(msg)
		return nil
	}
}
