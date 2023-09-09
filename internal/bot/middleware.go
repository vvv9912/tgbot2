package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/botkit"
)

// сюда кэш можно
func MwAdminOnly(next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

		next(ctx, bot, update)
		//проверка на админа
		return nil
	}
}

func MwUsersOnly(u botkit.UsersStorager, next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		//проверка на user
		uStatus, err := u.GetStatusUserByTgID(ctx, update.FromChat().ID)

		if err != nil {
			if err.Error() == NoRows {
				//ok
				err = nil
			} else {
				fmt.Println(err)
				return err
			}
		}
		switch uStatus {
		case UNoUser:
			//Добавить в бд и Приветствие!
		case UAdmin:

		case UUser:

		}
		next(ctx, bot, update)
		return nil
	}
}
