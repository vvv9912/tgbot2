package botkit

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/model"
	"tgbotv2/internal/storage"
)

//логика работы с ботом
import (
	"log"
	"runtime/debug"
	"time"
)

type OrderStorager interface {
	AddOrders(ctx context.Context, order model.Orders) error
	OrdersByTgID(ctx context.Context, tgId int64) ([]model.Orders, error)
}
type CorzinaStorager interface {
	AddCorzinas(ctx context.Context, corz model.Corzine) error
	CorzinaByTgId(ctx context.Context, tgId int64) ([]model.Corzine, error)
	UpdateCorzinaByTgId(ctx context.Context, tgId int64, article int, quantity int) error
	CorzinaByTgIdANDAtricle(ctx context.Context, tgId int64, article int) (model.Corzine, error)
	CorzinaByTgIdwithCalalog(ctx context.Context, tgId int64) ([]storage.DbCorzineCatalog, error)
	DeleteCorzinaByTgID(ctx context.Context, tgId int64) error
	DeleteCorzinaByTgIDandArticle(ctx context.Context, tgId int64, article int) error
}
type ProductsStorager interface {
	Catalog(ctx context.Context) ([]string, error)
	ProductsByCatalog(ctx context.Context, ctlg string) ([]model.Products, error)
	ProductByArticle(ctx context.Context, article int) (model.Products, error)
	AddProduct(ctx context.Context, product model.Products) error
}
type UsersStorager interface {
	GetStatusUserByTgID(ctx context.Context, tgID int64) (int, int, error)
	AddUser(ctx context.Context, users model.Users) error
	UpdateStateByTgID(ctx context.Context, tgId int64, state int) error
	//GetCorzinaByTgID(ctx context.Context, tgID int64) ([]int64, error)
	//UpdateCorzinaByTgId(ctx context.Context, tgId int64, corzina []int64) error
}
type BotCommand struct {
	Cmd  string `json:"cmd,omitempty"`
	Data string `json:"data,omitempty"` //в дата упоквано в зависимости от сообщения модель
}
type Bot struct {
	api           *tgbotapi.BotAPI
	cmdViews      map[string]ViewFunc // комманды тг бота
	textViews     map[string]ViewFunc
	callbackViews map[string]ViewFunc

	//PStorager     ProductsStorager
}
type BotInfo struct {
	TgId     int64
	IdStatus int
	IdState  int
}

// addcatalog
// listZakazov
// deleteCatalog id
type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo BotInfo) error

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{api: api}
}

//	type i interface {
//		RegisterTextView(string, ViewFunc)
//	}
//
//	func (b *Bot) Reegg(cmd string, ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update)) {
//		fmt.Println("aa")
//	}
func (b *Bot) RegisterTextView(cmd string, view ViewFunc) {
	if b.textViews == nil {
		b.textViews = make(map[string]ViewFunc)
	}
	b.textViews[cmd] = view
}
func (b *Bot) RegisterCallbackView(cmd string, view ViewFunc) {
	if b.callbackViews == nil {
		b.callbackViews = make(map[string]ViewFunc)
	}
	b.callbackViews[cmd] = view
}
func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}
	b.cmdViews[cmd] = view
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() { //перехват паники
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()
	//проверка авторизации пользователя
	//проверка из бд или кэша пользователя
	//
	var view ViewFunc
	if update.Message == nil {
		if update.CallbackQuery != nil {
			//тут чтото
			var Data BotCommand

			err := json.Unmarshal([]byte(update.CallbackQuery.Data), &Data)
			if err != nil {
				log.Printf("[ERROR] Json преобразование callback %v", err)
				return
			}
			clbck := Data.Cmd
			callbackView, ok := b.callbackViews[clbck]
			if !ok {
				return
			}
			view = callbackView
		} else if update.InlineQuery != nil {
			//тут  чтото
		} else {

			return
		}
	} else {
		if update.Message.IsCommand() {
			cmd := update.Message.Command()

			cmdView, ok := b.cmdViews[cmd]
			if !ok {
				return
			}

			view = cmdView

		} else if update.Message.Document != nil {

			cmd := update.Message.Caption
			cmdView, ok := b.cmdViews[cmd]
			if !ok {
				return
			}

			view = cmdView
		} else {
			//Если текстовая команда
			text := update.Message.Text
			if text == "" {
				return
			}
			textView, ok := b.textViews[text]
			if !ok {
				return
			}
			view = textView

		}
	}
	var botInfo BotInfo
	botInfo.TgId = update.FromChat().ID
	if err := view(ctx, b.api, update, botInfo); err != nil {
		log.Printf("[ERROR] failed to handle update: %v", err)

		b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintln(err)))
		if _, err := b.api.Send(
			tgbotapi.NewMessage(update.Message.Chat.ID, "internal error"),
		); err != nil {
			log.Printf("[ERROR] failed to send message: %v", err)
		}
	}
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	//https://yandex.ru/search/?text=webhook+telegram+bot+api+golang&lr=213&clid=2574587&win=575 webhook
	for {
		select {
		case update := <-updates:
			go func(update tgbotapi.Update) {
				updateCtx, updateCancel := context.WithTimeout(ctx, 60*time.Second)
				b.handleUpdate(updateCtx, update)
				//select c таймаутом. То что запрос долго обрабатывается
				updateCancel()
			}(update)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
