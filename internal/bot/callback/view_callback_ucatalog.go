package callback

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tgbotv2/internal/botkit"
)

func ViewCallbackUcatalog(s botkit.ProductsStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
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

			//data = fmt.Sprintf("/ucatalog\ncategory:%s", catalog[i]) //надо передавать команду +id что удалить(?) //todo передать сразу json

			if err != nil {
				log.Println("") //todo
			}

			if len(sProducts[i].PhotoUrl) != 0 {

				//todo photo
				//var mediaGroup []interface{}
				//for k := range sProducts[i].PhotoUrl {
				var mediaGroup []interface{}
				//if len(sProducts[i].PhotoUrl) > 1 {
				for k := range sProducts[i].PhotoUrl {
					photoFileBytes := tgbotapi.FileBytes{
						Name:  "photo1.jpg",
						Bytes: sProducts[i].PhotoUrl[k],
					}
					photo1 := tgbotapi.NewInputMediaPhoto(photoFileBytes)
					if k == 0 {
						photo1.Caption = sProducts[i].Name //описание строго под одной фото
					}

					mediaGroup = append(mediaGroup, photo1)
				}

				msg := tgbotapi.NewMediaGroup(update.CallbackQuery.From.ID, mediaGroup)
				bot.Send(msg)
				//} else {
				//	photoFileBytes := tgbotapi.FileBytes{
				//		Name:  "photo2.jpg",
				//		Bytes: sProducts[i].PhotoUrl[0],
				//	}
				//	//photo1 := tgbotapi.NewInputMediaPhoto(photoFileBytes)
				//	msg := tgbotapi.NewPhoto(update.CallbackQuery.From.ID, photoFileBytes)
				//	bot.Send(msg)
				//}
				//bot.Send(msg) //json: cannot unmarshal array into Go value of type tgbotapi.Message

				dataAddCorz := AddCorzine{
					Article: sProducts[i].Article,
				}
				msgAddCorz, err := json.Marshal(dataAddCorz)
				if err != nil {
					log.Println("") //todo
				}
				data := botkit.BotCommand{Cmd: "/moredetailed",
					Data: string(msgAddCorz)}
				podrobnee, err := json.Marshal(data)
				dataMsg := botkit.BotCommand{
					Cmd:  "/addCorzine",
					Data: string(msgAddCorz), //не над
				}
				sss, err := json.Marshal(dataMsg)
				if err != nil {
					log.Println("") //todo
				}

				var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Подробнее", string(podrobnee)),
						tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", string(sss)),
					),
				)
				msg2 := tgbotapi.NewMessage(update.CallbackQuery.From.ID, text)

				_ = numericKeyboardInline
				msg2.ReplyMarkup = numericKeyboardInline
				_, err = bot.Send(msg2) //todo
				if err != nil {
					log.Println(err)
					return err
				}
				//	}
			} else { //если нет фото
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
				//text = `<a href="tg://btn/dmVyeSBsb25nIHN0cmluZwBldmVuIGxvbmdlciBzdHJpbmcAZXZlsadssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssbiBsb25nZXIgc3RpbGwgc3RyaW5n">\u200b</a>` + text
				//msg.ParseMode = "HTML"
				//sss := fmt.Sprintf("/addCorzine\narticle:%d\ncategory:%s", sProducts[i].Article, sProducts[i].Catalog) //todo
				//Завернули структуру в структуру (джсон в джсон)
				dataAddCorz := AddCorzine{
					Article: sProducts[i].Article,
				}
				msgAddCorz, err := json.Marshal(dataAddCorz)
				if err != nil {
					log.Println("") //todo
				}
				dataMsg := botkit.BotCommand{
					Cmd:  "/addCorzine",
					Data: string(msgAddCorz),
				}
				data := botkit.BotCommand{Cmd: "/moredetailed",
					Data: string(msgAddCorz)}
				podrobnee, err := json.Marshal(data)
				//_ = msgAddCorz
				//dataMsg := botkit.BotCommand{
				//	Cmd: "/addCorzine",
				//}
				sss, err := json.Marshal(dataMsg)

				var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Подробнее", string(podrobnee)),
						tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", string(sss)),
					),
				)
				log.Println("len(sss)=", len(sss))
				msg.ReplyMarkup = numericKeyboardInline
				_, err = bot.Send(msg) //todo
				if err != nil {
					log.Println(err)
					return err
				}

			}
		}
		return nil
	}
}
