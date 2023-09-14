package cmd

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mholt/archiver"
	"io"
	"net/http"
	"os"
	"path"
	"tgbotv2/internal/botkit"
	"tgbotv2/internal/config"
)

// Скачивание файла
func ViewCmdAdddbfile() botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		//Проверка типа, что это rar
		fileid := update.Message.Document.FileID
		var file tgbotapi.FileConfig
		file.FileID = fileid
		File, err := bot.GetFile(file)

		if err != nil {
			fmt.Println(err)
			return err
		}
		//File.FilePath
		//out, err := os.Create(File)
		//defer out.Close()
		//n, err := io.Copy(out, resp.Body)
		pathserv := "https://api.telegram.org/file/" + "bot" + config.Get().TelegramBotToken + "/" + File.FilePath
		resp, err := http.Get(pathserv)
		_, err = os.Stat("tmp_downl")
		if !os.IsNotExist(err) {

			err = os.RemoveAll("tmp_downl")
			if err != nil {
				fmt.Println(err)
				return err
			}
			//err = os.MkdirAll("tmp_downl", os.ModePerm)
			//if err != nil {
			//	fmt.Println(err)
			//	return err
			//}
		}
		err = os.Mkdir("tmp_downl", os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
		FileName := path.Join("tmp_downl", "db.zip")
		//	FileName := "tmp_downl/db.RAR"

		out, err := os.Create(FileName)
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}
		//a, err := unarr.NewArchive(FileName)
		//if err != nil {
		//	fmt.Println(err)
		//	return err
		//}
		//_, err = a.Extract("tmp_downl/")
		//if err != nil {
		//	fmt.Println(err)
		//	return err
		//}
		err = archiver.Unarchive(FileName, "./tmp_downl/")
		// path.Dir("tmp/")
		if err != nil {
			fmt.Println(err)
			return err
		}

		// сюда протащить загрузку в бд!
		return nil
	}
}
