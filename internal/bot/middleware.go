package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tgbotv2/internal/botkit"
	"tgbotv2/internal/model"
)

// сюда кэш можно
func MwAdminOnly(next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {

		next(ctx, bot, update, botInfo)
		//проверка на админа
		return nil
	}
}

func MwUsersOnly(u botkit.UsersStorager, next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		//проверка на user
		uStatus, uState, err := u.GetStatusUserByTgID(ctx, update.FromChat().ID)
		if err != nil {
			if err.Error() == NoRows {
				//ok
				err = nil
			} else {
				fmt.Println(err)
				return err
			}
		}
		botInfo.IdStatus = uStatus
		botInfo.IdState = uState

		switch uStatus {
		case UNoUser:
			err = u.AddUser(ctx, model.Users{
				TgID:       botInfo.TgId,
				StateUser:  UUser,
				StatusUser: E_STATE_NOTHING,
			})
			if err != nil {
				log.Println(err)
				return err
			}
			botInfo.IdState = UUser
			botInfo.IdStatus = E_STATE_NOTHING
			//Добавить в бд и Приветствие!
		case UAdmin:
			//
		case UUser:
			//
		}
		next(ctx, bot, update, botInfo)
		return nil
	}
}
