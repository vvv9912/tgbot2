package main

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
	bot "tgbotv2/internal/bot/callback"
	callback "tgbotv2/internal/bot/callback/adddb"
	"tgbotv2/internal/bot/cmd"
	"tgbotv2/internal/bot/middleware"
	"tgbotv2/internal/bot/text"
	"tgbotv2/internal/botkit"
	"tgbotv2/internal/config"
	"tgbotv2/internal/storage"
)

func main() {

	//rr := exel.Read()
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	//_ = rr
	if err != nil {
		log.Printf("failed to create bot:%v", err)
		return
	}
	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("failed to connect db:%v", err)
		return
	}
	defer db.Close()
	sProducts := storage.NewProductsPostgresStorage(db)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	sUsers := storage.NewUsersPostgresStorage(db)
	sCorzina := storage.NewCorzinaPostgresStorage(db)
	sOrder := storage.NewOrdersPostgresStorage(db)
	Ourbot := botkit.New(botAPI)
	mw := middleware.NewMiddleware(sUsers)
	//
	Ourbot.RegisterCmdView("start", mw.MwUsersOnly(cmd.ViewCmdStart(cmd.ViewCmdButton())))
	Ourbot.RegisterCmdView("button", cmd.ViewCmdButton())
	Ourbot.RegisterCmdView("adminbutton", cmd.ViewCmdAdminButton())
	Ourbot.RegisterCmdView("database", cmd.ViewCmdAddDatabase())          //Проверка что админ
	Ourbot.RegisterCmdView("/adddbfile", cmd.ViewCmdAdddbfile(sProducts)) // Проверка что админ
	//
	Ourbot.RegisterTextView("Каталог", mw.MwUsersOnly(text.ViewTextCatalog(sProducts)))
	Ourbot.RegisterTextView("Корзина", mw.MwUsersOnly(text.ViewTextCorzine(sCorzina)))
	//
	Ourbot.RegisterCallbackView("/ucatalog", bot.ViewCallbackUcatalog(sProducts))
	Ourbot.RegisterCallbackView("/addCorzine", bot.ViewCallbackAddcorzine(sCorzina))
	Ourbot.RegisterCallbackView("/moredetailed", bot.ViewCallbackMoredetailed(sProducts, sCorzina))
	Ourbot.RegisterCallbackView("/adddatabase", callback.ViewCallbackAdddatabase(sUsers)) //проверка что админ
	Ourbot.RegisterCallbackView("/deleteshoppcart", bot.ViewCallbackdeleteshoppcart(sCorzina))
	Ourbot.RegisterCallbackView(" /addorder", bot.ViewCallbackAddorder(sOrder, sCorzina))
	//deleteshoppcartposition
	Ourbot.RegisterCallbackView("/deleteshoppcartposition", bot.ViewCallbackdeleteshoppcartposition(sCorzina))
	//
	if err := Ourbot.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] failed to start bot: %v", err)
			return
		}
		log.Println("bot stopped")
	}
}
