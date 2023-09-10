package callback

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tgbotv2/internal/botkit"
)

func ViewCallbackAddcorzine(s botkit.ProductsStorager, u botkit.UsersStorager, c botkit.CorzinaStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {

		var Data botkit.BotCommand
		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &Data)
		if err != nil {
			log.Printf("[ERROR] Json преобразование callback %v", err)
			return err
		}
		var MsgAddCorzine AddCorzine
		err = json.Unmarshal([]byte(Data.Data), &MsgAddCorzine)
		if err != nil {
			log.Printf("[ERROR] Json преобразование callback %v", err)
			return err
		}
		//Добавление в БД
		//u.GetCorzinaByTgID()
		//u.UpdateCorzinaByTgId()
		return nil
	}
}
