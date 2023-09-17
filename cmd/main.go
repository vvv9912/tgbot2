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
	"tgbotv2/internal/bot/callback"

	callback2 "tgbotv2/internal/bot/callback/adddb"
	callback3 "tgbotv2/internal/bot/callback/order"

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
	Ourbot.RegisterCmdView("button", cmd.ViewCmdButton()) //check user
	Ourbot.RegisterCmdView("adminbutton", cmd.ViewCmdAdminButton())
	Ourbot.RegisterCmdView("database", cmd.ViewCmdAddDatabase())          //Проверка что админ
	Ourbot.RegisterCmdView("/adddbfile", cmd.ViewCmdAdddbfile(sProducts)) // Проверка что админ
	//
	Ourbot.RegisterTextView("Каталог", mw.MwUsersOnly(text.ViewTextCatalog(sProducts)))
	Ourbot.RegisterTextView("Корзина", mw.MwUsersOnly(text.ViewTextCorzine(sCorzina)))
	Ourbot.RegisterTextView("Мои заказы", mw.MwUsersOnly(text.ViewTextOrder(sOrder)))
	//
	Ourbot.RegisterCallbackView("/ucatalog", callback.ViewCallbackUcatalog(sProducts))                   //check user
	Ourbot.RegisterCallbackView("/addCorzine", callback.ViewCallbackAddcorzine(sCorzina))                //check user
	Ourbot.RegisterCallbackView("/moredetailed", callback.ViewCallbackMoredetailed(sProducts, sCorzina)) //check user
	Ourbot.RegisterCallbackView("/adddatabase", callback2.ViewCallbackAdddatabase(sUsers))               //проверка что админ
	Ourbot.RegisterCallbackView("/deleteshoppcart", callback.ViewCallbackdeleteshoppcart(sCorzina))      //check user
	//deleteshoppcartpositionViewCallbackAddorder(sOrder, sCorzina))
	Ourbot.RegisterCallbackView("/deleteshoppcartposition", callback.ViewCallbackdeleteshoppcartposition(sCorzina)) //check user

	Ourbot.RegisterCallbackView("/addorder", callback3.ViewCallbackAddorder())                                   //check user
	Ourbot.RegisterCallbackView("/addOrderSamovivoz", callback3.ViewCallbackAddOrderSamovivoz(sOrder, sCorzina)) //check user

	//
	Ourbot.RegisterTextView("FAQ", text.ViewTextFaq())
	if err := Ourbot.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] failed to start bot: %v", err)
			return
		}
		log.Println("bot stopped")
	}
}
