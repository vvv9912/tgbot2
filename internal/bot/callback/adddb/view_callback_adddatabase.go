package callback_adddb

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/bot/constant"
	"tgbotv2/internal/botkit"
)

func ViewCallbackAdddatabase(u botkit.UsersStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		//Тут добавляем нашему админу состояние.
		//состояние 100 - ждем файл

		err := u.UpdateStateByTgID(ctx, botInfo.TgId, constant.U_STATE_ADDDB_WAIT_FILE)
		if err != nil {
			fmt.Println("[ERROR] Update state: ", err)
			return err
		}
		txt := "Добавьте файл с базой данных в соответсвутющем формате (см. инструкцию). Обязательно напишите комнаду. /adddbfile"
		msg := tgbotapi.NewMessage(botInfo.TgId, txt)
		bot.Send(msg)
		//structMsgAddDb := botkit.BotCommand{
		//	Cmd: "/adddatabase",
		//}
		//msgAddDb, err := json.Marshal(structMsgAddDb)
		//if err != nil {
		//	log.Println("") //todo
		//}
		//var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
		//	tgbotapi.NewInlineKeyboardRow(
		//		tgbotapi.NewInlineKeyboardButtonData("Добавить бд", string(msgAddDb)),
		//		//tgbotapi.NewInlineKeyboardButtonData("Добавить бд", string(sss)),
		//	),
		//)
		//msg := tgbotapi.NewMessage(update.FromChat().ID, "Изменение базы данных")
		//msg.ReplyMarkup = numericKeyboardInline
		//if _, err := bot.Send(msg); err != nil {
		//	return err
		//}
		return nil
	}
}
