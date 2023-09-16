package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/botkit"
	"tgbotv2/internal/model"
	"time"
)

func ViewCallbackAddorder(o botkit.OrderStorager, c botkit.CorzinaStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		//выбрать доставку https://github.com/vvv9912/mytgbot/blob/main/tg2/parscommand.go#L714
		//написать адм или выбрать пвз и тп...
		//потом
		corz, err := c.CorzinaByTgId(ctx, botInfo.TgId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		//
		_ = corz
		err = o.AddOrders(ctx, model.Orders{
			ID:          0,
			TgID:        0,
			StatusOrder: 0,
			Pvz:         "",
			Order:       "",
			CreatedAt:   time.Now().UTC().Add(3 * time.Hour),
			ReadAt:      time.Now().UTC().Add(3 * time.Hour),
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}
}
