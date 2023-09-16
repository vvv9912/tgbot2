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

func ViewCallbackdeleteshoppcartposition(c botkit.CorzinaStorager) botkit.ViewFunc {
	//view_callback_deleteshoppcartposition.go
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

		corz, err := c.CorzinaByTgIdANDAtricle(ctx, botInfo.TgId, article)
		if err != nil {
			fmt.Println("CorzinaByTgIdANDAtricle :", err)
			return err
		}
		if corz == (model.Corzine{}) {
			fmt.Println("corz == (model.Corzine{}) :", err)
			return err
		}
		if corz.Quantity == 1 {
			err := c.DeleteCorzinaByTgIDandArticle(ctx, botInfo.TgId, article)
			if err != nil {
				fmt.Println("DeleteCorzinaByTgIDandArticle :", err)
				return err
			}
		} else {
			newQuan := corz.Quantity - 1
			err := c.UpdateCorzinaByTgId(ctx, botInfo.TgId, article, newQuan)
			if err != nil {
				fmt.Println("UpdateCorzinaByTgId :", err)
				return err
			}
		}
		//err := c.DeleteCorzinaByTgID(ctx, botInfo.TgId)
		//if err != nil {
		//	fmt.Println("DeleteCorzinaByTgID :", err)
		//	return err
		//}
		answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Позиция удалена!"}
		bot.Send(answ)
		return nil
	}
}
