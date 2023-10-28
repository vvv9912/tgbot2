package file

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"tgbotv2/internal/config"
)

func savefile(outfile tgbotapi.File) (string, error) {
	nameid := uuid.New()
	pathserv := "https://api.telegram.org/file/" + "bot" + config.Get().TelegramBotToken + "/" + outfile.FilePath
	w := strings.Split(outfile.FilePath, ".") //todo
	fileext := w[len(w)-1]
	//nameid + "." + fileext
	filename := fmt.Sprintf("%s.%s", &nameid, fileext)
	p1 := "database"
	p2 := "photo"
	f := path.Join(p1, p2, filename)
	resp, err := http.Get(pathserv)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	out, err := os.Create(f)
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	return f, err
}
