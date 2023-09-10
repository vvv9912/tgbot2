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
	"tgbotv2/internal/bot/cmd"
	"tgbotv2/internal/bot/middleware"
	"tgbotv2/internal/bot/text"
	"tgbotv2/internal/botkit"
	"tgbotv2/internal/config"
	"tgbotv2/internal/storage"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
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
	Ourbot := botkit.New(botAPI)
	mw := middleware.NewMiddleware(sUsers)
	Ourbot.RegisterCmdView("start", mw.MwUsersOnly(cmd.ViewCmdStart(cmd.ViewCmdButton())))
	Ourbot.RegisterCmdView("button", cmd.ViewCmdButton())
	Ourbot.RegisterCmdView("adminbutton", cmd.ViewCmdAdminButton())

	Ourbot.RegisterTextView("Каталог", mw.MwUsersOnly(text.ViewTextCatalog(sProducts)))

	Ourbot.RegisterCallbackView("/ucatalog", bot.ViewCallbackUcatalog(sProducts))

	if err := Ourbot.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] failed to start bot: %v", err)
			return
		}
		log.Println("bot stopped")
	}
}
