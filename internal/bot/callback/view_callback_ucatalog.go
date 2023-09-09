package bot

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"strconv"
	"tgbotv2/internal/botkit"
)

func ViewCallbackUcatalog(s botkit.ProductsStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		var Data botkit.BotCommand
		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &Data)
		if err != nil {
			log.Printf("[ERROR] Json преобразование callback %v", err)
			return err
		}
		//Тут запрос в БД и поиск артикулей по названию каталога
		sProducts, err := s.ProductsByCatalog(ctx, Data.Data)
		if err != nil {
			log.Printf("[ERROR] Get products from ProductsStorage %v", err)

			return err
		}
		if len(sProducts) == 0 {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Товары отсутсвуют!")
			if _, err := bot.Send(msg); err != nil {
				return err
			}
		}
		//
		for i := range sProducts {
			//text := fmt.Sprintf("Артикул: %d\nНазвание: %s\n%s\nЦена: %0.2fрублей\n", db_send[i].Article, db_send[i].Name, db_send[i].Description, db_send[i].Price)
			text := fmt.Sprintf("Артикул: %d\nНазвание: %s\nПодходит для: \nЦена: %0.2f рублей\n", sProducts[i].Article, sProducts[i].Name, sProducts[i].Price)
			//podrobnee := fmt.Sprintf("/moredetailed\narticle:%d", sProducts[i].Article) //todo toJson
			data := botkit.BotCommand{Cmd: "/moredetailed",
				Data: strconv.Itoa(sProducts[i].Article)}
			//data = fmt.Sprintf("/ucatalog\ncategory:%s", catalog[i]) //надо передавать команду +id что удалить(?) //todo передать сразу json
			podrobnee, err := json.Marshal(data)
			if err != nil {
				log.Println("") //todo
			}

			if sProducts[i].PhotoUrl != "" {
				photoBytes, err := ioutil.ReadFile(sProducts[i].PhotoUrl)
				_ = err
				photoFileBytes := tgbotapi.FileBytes{
					Name:  "photo",
					Bytes: photoBytes,
				}

				msg := tgbotapi.NewPhoto(update.CallbackQuery.From.ID, nil)
				msg.File = photoFileBytes

				//sss := fmt.Sprintf("/addCorzine\narticle:%d\ncategory:%s", sProducts[i].Article, sProducts[i].Catalog) //todo
				//Завернули структуру в структуру (джсон в джсон)
				dataAddCorz := AddCorzine{
					Article:  sProducts[i].Article,
					Category: sProducts[i].Catalog,
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
				msg.Caption = text
				var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Подробнее", string(podrobnee)),
						tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", string(sss)),
					),
				)

				msg.ReplyMarkup = numericKeyboardInline
				bot.Send(msg) //todo
			} else { //если нет фото
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
				//sss := fmt.Sprintf("/addCorzine\narticle:%d\ncategory:%s", sProducts[i].Article, sProducts[i].Catalog) //todo
				//Завернули структуру в структуру (джсон в джсон)
				dataAddCorz := AddCorzine{
					Article:  sProducts[i].Article,
					Category: sProducts[i].Catalog,
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
						tgbotapi.NewInlineKeyboardButtonData("Подробнее", string(podrobnee)),
						tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", string(sss)),
					),
				)
				msg.ReplyMarkup = numericKeyboardInline
				bot.Send(msg) //todo
			}
		}
		return nil
	}
}
