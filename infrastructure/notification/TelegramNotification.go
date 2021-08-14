package notification

import (
	"idealista/domain/ports"
	"os"
	"path"
	"strconv"

	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramNotification struct {
}

func NewTelegramNotification() ports.NotificationService {
	return &telegramNotification{}
}

func (tn telegramNotification) SendReports() {
	telegramChatId, bot := getTelegramCredentials()

	fileDir, _ := os.Getwd()
	tn.sendFile(ports.RENTAL_REPORT_MONTHLY, bot, fileDir, telegramChatId)
	tn.sendFile(ports.SALE_REPORT_MONTHLY, bot, fileDir, telegramChatId)
}

func getTelegramCredentials() (int64, *tgbotapi.BotAPI) {
	telegramToken := os.Getenv("TELEGRAM_API_TOKEN")
	telegramChatId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)

	if telegramToken == "" || telegramChatId == 0 {
		zap.L().Panic("Got empty telegram credentials")
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	bot.Debug = true

	if err != nil {
		zap.L().Panic("Could not get telegram API instance", zap.Error(err))
	}

	return telegramChatId, bot
}

func (tn telegramNotification) sendFile(fileName string, bot *tgbotapi.BotAPI, fileDir string, chatId int64) {
	filePath := path.Join(fileDir, fileName)

	msg := tgbotapi.NewPhoto(chatId, filePath)
	_, err2 := bot.Send(msg)

	if err2 != nil {
		zap.L().Error("Document not sent", zap.Error(err2))
	}
}
