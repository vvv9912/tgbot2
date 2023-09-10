package middleware

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tgbotv2/internal/bot/constant"

	"tgbotv2/internal/botkit"
	"tgbotv2/internal/model"
)

type Middleware struct {
	UserStorage botkit.UsersStorager
}

func NewMiddleware(userStorage botkit.UsersStorager) *Middleware {
	return &Middleware{UserStorage: userStorage}
}

// сюда кэш можно
func (m *Middleware) MwAdminOnly(next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {

		next(ctx, bot, update, botInfo)
		//проверка на админа
		return nil
	}
}

func (m *Middleware) MwUsersOnly(next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		//проверка на user

		uStatus, uState, err := m.UserStorage.GetStatusUserByTgID(ctx, update.FromChat().ID)
		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
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
		case constant.UNoUser:
			err = m.UserStorage.AddUser(ctx, model.Users{
				TgID:       botInfo.TgId,
				StateUser:  constant.UUser,
				StatusUser: constant.E_STATE_NOTHING,
			})
			if err != nil {
				log.Println(err)
				return err
			}
			botInfo.IdState = constant.UUser
			botInfo.IdStatus = constant.E_STATE_NOTHING
			//Добавить в бд и Приветствие!
		case constant.UAdmin:
			//
		case constant.UUser:
			//
		}

		next(ctx, bot, update, botInfo)

		return nil
	}
}
