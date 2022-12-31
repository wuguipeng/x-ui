package service

import (
	"fmt"
	"x-ui/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	settingService SettingService
	serverService  ServerService
	lastStatus     *Status
}

var bot *tgbotapi.BotAPI

func (j *TelegramService) Init() {
	//Telegram bot basic info
	tgBottoken, err := j.settingService.GetTgBotToken()
	if err != nil {
		logger.Warning("sendMsgToTgbot failed,GetTgBotToken fail:", err)
	}
	// tgBotid, err := j.settingService.GetTgBotChatId()
	// if err != nil {
	// 	logger.Warning("sendMsgToTgbot failed,GetTgBotChatId fail:", err)
	// 	return
	// }

	botInit, err := tgbotapi.NewBotAPI(tgBottoken)
	if err != nil {
		fmt.Println("get tgbot error:", err)
	}
	botInit.Debug = true
	bot = botInit
}

func (j *TelegramService) Command() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		// video := update.Message.Photo
		// fmt.Print(video)
		// photoSize := video[len(video)-1]
		// fileConfig := tgbotapi.FileConfig{
		// 	FileID: photoSize.FileID,
		// }
		// video := update.Message.Video
		// fileConfig := tgbotapi.FileConfig{
		// 	FileID: video.FileID,
		// }
		// file, _ := bot.GetFile(fileConfig)
		// // filePath := file.FilePath
		// fmt.Printf(file.FilePath)
		// fmt.Printf(file.Link("5471229987:AAFhn7cwh-KX2brKtF868z4TdLfmrcISmGY"))

		// file, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)

		// ddd, _ := base64.StdEncoding.DecodeString(resp.Result.MarshalJSON()) //成图片文件并把文件写入到buffer
		// err2 := ioutil.WriteFile("./output.jpg", ddd, 0666)

		// if update.Message.IsCommand() {
		// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		// 	switch update.Message.Command() {
		// 	case "help":
		// 		msg.Text = "type /status."
		// 	case "status":
		// 		status := j.serverService.GetStatus(j.lastStatus)
		// 		text := fmt.Sprintf("CPU:%s\r\n", status.Cpu)
		// 		// text := fmt.Sprintf("")
		// 		msg.Text = text
		// 	case "withArgument":
		// 		msg.Text = "You supplied the following argument: " + update.Message.CommandArguments()
		// 	case "html":
		// 		msg.ParseMode = "html"
		// 		msg.Text = "This will be interpreted as HTML, click <a href=\"https://www.example.com\">here</a>"
		// 	default:
		// 		msg.Text = "I don't know that command"
		// 	}
		// 	bot.Send(msg)
		// }
	}
}

func (j *TelegramService) RegisterCommand() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	command1 := tgbotapi.BotCommand{
		Command:     "/help",
		Description: "帮助",
	}
	command2 := tgbotapi.BotCommand{
		Command:     "/status",
		Description: "获取状态",
	}
	myCommand := tgbotapi.NewSetMyCommands(command1, command2)
	bot.Request(myCommand)
}
