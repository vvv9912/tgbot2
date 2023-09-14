package callback

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"tgbotv2/internal/botkit"
	"tgbotv2/internal/model"
)

func ViewCallbackMoredetailed(p botkit.ProductsStorager, c botkit.CorzinaStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {

		var Data botkit.BotCommand
		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &Data)
		if err != nil {
			log.Printf("[ERROR] Json преобразование callback %v", err)
			return err
		}
		article, err := strconv.Atoi(Data.Data)
		if err != nil {
			log.Println("[ERROR] Strconv in MoreDetailed %v", err)
			return err
		}
		product, err := p.ProductByArticle(ctx, article)
		//Обновлять предыдущее сообщение (добавлять подробности)
		if (product == model.Products{}) {
			return err //todo
		}
		if product.PhotoUrl != "" {

			text := fmt.Sprintf("Артикул: %d\nНазвание: %s\n%s\nЦена: %0.2fрублей\n", product.Article, product.Name, product.Description, product.Price)
			ms1 := tgbotapi.NewEditMessageCaption(botInfo.TgId, update.CallbackQuery.Message.MessageID, text)
			dataAddCorz := AddCorzine{
				Article: product.Article,
			}
			msgAddCorz, err := json.Marshal(dataAddCorz)
			if err != nil {
				log.Println("") //todo
			}
			dataMsg := botkit.BotCommand{
				Cmd:  "/addCorzine",
				Data: string(msgAddCorz),
			}
			sss, err := json.Marshal(dataMsg)
			var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", string(sss)),
				),
			)
			ms1.ReplyMarkup = &numericKeyboardInline
			_, err = bot.Send(ms1)
			if err != nil {
				log.Println("[ERROR] Send bot: %v", err) //todo
				return err
			}
		} else {
			text := fmt.Sprintf("Артикул: %d\nНазвание: %s\n%s\nЦена: %0.2fрублей\n", product.Article, product.Name, product.Description, product.Price)
			ms1 := tgbotapi.NewEditMessageText(botInfo.TgId, update.CallbackQuery.Message.MessageID, text)
			//	ms1 := tgbotapi.NewEditMessageText(int64(tg_id), update.CallbackQuery.Message.MessageID, text)
			dataAddCorz := AddCorzine{
				Article: product.Article,
			}
			msgAddCorz, err := json.Marshal(dataAddCorz)
			if err != nil {
				log.Println("") //todo
				return err
			}
			dataMsg := botkit.BotCommand{
				Cmd:  "/addCorzine",
				Data: string(msgAddCorz),
			}
			sss, err := json.Marshal(dataMsg)
			var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", string(sss)),
				),
			)
			ms1.ReplyMarkup = &numericKeyboardInline
			_, err = bot.Send(ms1)
			if err != nil {
				log.Println("[ERROR] Send bot: %v", err) //todo
				return err
			}
		}

		return nil
	}
}
