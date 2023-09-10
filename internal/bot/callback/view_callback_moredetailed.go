package callback

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tgbotv2/internal/botkit"
	"tgbotv2/internal/model"
	"time"
)

func ViewCallbackMoredetailed(c botkit.CorzinaStorager) botkit.ViewFunc {

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
		corz, err := c.CorzinaByTgIdANDAtricle(ctx, botInfo.TgId, MsgAddCorzine.Article)
		_ = corz

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = nil
				err = c.AddCorzinas(ctx, model.Corzine{
					TgId:      botInfo.TgId,
					Article:   MsgAddCorzine.Article,
					Quantity:  1,
					CreatedAt: time.Time{},
				})
				if err != nil {
					log.Println(err)
					return err
				}
				//send message 'add corzine' todo
				msg := tgbotapi.NewMessage(botInfo.TgId, "Добавлено в корзину!")
				_, err = bot.Send(msg)
				if err != nil {
					fmt.Println(err)
					return err
				}
				answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Добавлено в корзину!"}
				_, err = bot.Send(answ)
				if err != nil {
					fmt.Println(err)
					return err
				}
				return nil
			}
			log.Fatalln(err)
			return err
		}
		if corz == (model.Corzine{}) { //Проверка на всякий случай
			err = c.AddCorzinas(ctx, model.Corzine{
				TgId:      botInfo.TgId,
				Article:   MsgAddCorzine.Article,
				Quantity:  1,
				CreatedAt: time.Now().UTC().Add(3 * time.Hour),
			})
			if err != nil {
				log.Println(err)
				return err
			}
			msg := tgbotapi.NewMessage(botInfo.TgId, "Добавлено в корзину!")
			_, err = bot.Send(msg)
			if err != nil {
				fmt.Println(err)
				return err
			}
			answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Добавлено в корзину!"}
			_, err = bot.Send(answ)
			if err != nil {
				fmt.Println(err)
				return err
			}
			//send message 'add corzine' todo
			return nil
		}
		corz.Quantity++
		corz.CreatedAt = time.Now().UTC().Add(3 * time.Hour)
		err = c.UpdateCorzinaByTgId(ctx, botInfo.TgId, corz.Article, corz.Quantity) //todo обновлять время
		if err != nil {
			log.Println(err)
			return err
		}
		msg := tgbotapi.NewMessage(botInfo.TgId, "Добавлено в корзину!")
		_, err = bot.Send(msg)
		if err != nil {
			fmt.Println(err)
			return err
		}
		answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Добавлено в корзину!"}
		_, err = bot.Send(answ)
		if err != nil {
			fmt.Println(err)
			return err
		}
		//send message 'add corzine' todo
		return nil
	}
}
