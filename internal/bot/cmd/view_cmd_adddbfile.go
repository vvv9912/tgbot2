package cmd

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mholt/archiver"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"tgbotv2/internal/botkit"
	"tgbotv2/internal/config"
	"tgbotv2/internal/exel"
	"tgbotv2/internal/model"
)

// Скачивание файла
func ViewCmdAdddbfile(p botkit.ProductsStorager) botkit.ViewFunc {

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
		//База данных.xlsx
		//"./tmp_downl/test_bd/База данных.xlsx"
		pathfiledb := path.Join("tmp_downl", "test_bd", "База данных.xlsx") //после архивации
		pathdir := path.Join("tmp_downl", "test_bd")
		ex := exel.NewExcel(pathfiledb)
		indb := ex.Read()
		if indb == nil {
			return errors.New("pars file db")
		}
		for i := range *indb {
			var photoBytes []byte
			if (*indb)[i].PhotoUrl != "" {
				//ппроверка неправильного пути
				if strings.Contains((*indb)[i].PhotoUrl, fmt.Sprintf(`\`)) {
					a := strings.Split((*indb)[i].PhotoUrl, `\`)
					var photoUrl string
					for i := range a {
						photoUrl = path.Join(photoUrl, a[i])
					}
					photoBytes, err = ioutil.ReadFile(path.Join(pathdir, photoUrl))
					if err != nil {
						return err
					}
				} else {
					photoBytes, err = ioutil.ReadFile(path.Join(pathdir, (*indb)[i].PhotoUrl))
					if err != nil {
						return err
					}
				}
			}
			err = p.AddProduct(ctx, model.Products{
				Article:     (*indb)[i].Article,
				Catalog:     (*indb)[i].Catalog,
				Name:        (*indb)[i].Name,
				Description: (*indb)[i].Description,
				PhotoUrl:    photoBytes,
				Price:       (*indb)[i].Price,
				Length:      (*indb)[i].Length,
				Width:       (*indb)[i].Width,
				Height:      (*indb)[i].Height,
				Weight:      (*indb)[i].Weight,
			})
			if err != nil {
				return err
			}
		}
		// 1. проверка на существование такого артикула
		//
		return nil
	}
}
