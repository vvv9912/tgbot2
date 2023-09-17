package callback_order

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tgbotv2/internal/botkit"
)

func ViewCallbackAddorder() botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		//выбрать доставку https://github.com/vvv9912/mytgbot/blob/main/tg2/parscommand.go#L714
		//написать адм или выбрать пвз и тп...
		//потом
		dataMsgSam := botkit.BotCommand{
			Cmd: "/addOrderSamovivoz",
			//Data: "",
		}
		msgSam, err := json.Marshal(dataMsgSam)
		if err != nil {
			log.Println("") //todo
			return err
		}
		dataMsgSdek := botkit.BotCommand{
			Cmd:  "/addOrderSdek",
			Data: "",
		}
		msgSdek, err := json.Marshal(dataMsgSdek)
		if err != nil {
			log.Println("") //todo
			return err
		}
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите способ доставки")
		var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Самовывоз", string(msgSam)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Доставка СДЭК", string(msgSdek)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Оформить заказ, о доставке с вами свяжутся(?)", string("s")),
			),
		)
		msg.ReplyMarkup = numericKeyboardInline
		bot.Send(msg)
		//corz, err := c.CorzinaByTgId(ctx, botInfo.TgId)
		//if err != nil {
		//	fmt.Println(err)
		//	return err
		//}
		////
		//_ = corz
		//err = o.AddOrders(ctx, model.Orders{
		//	ID:          0,
		//	TgID:        0,
		//	StatusOrder: 0,
		//	Pvz:         "",
		//	Order:       "",
		//	CreatedAt:   time.Now().UTC().Add(3 * time.Hour),
		//	ReadAt:      time.Now().UTC().Add(3 * time.Hour),
		//})
		//if err != nil {
		//	fmt.Println(err)
		//	return err
		//}
		return nil
	}
}
