package telegrambot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Open the Fucking Door!", "open"),
	),
)

// TelegramBot ...
type TelegramBot struct {
	config *Config
	logger *logrus.Logger
}

// New ...
func New(config *Config) *TelegramBot {
	return &TelegramBot{
		config: config,
		logger: logrus.New(),
	}
}

// Start ...
func (t *TelegramBot) Start() error {
	if err := t.configureLogger(); err != nil {
		return err
	}
	t.logger.Info("Starting Bot")

	t.runBot()

	return nil
}

func (t *TelegramBot) configureLogger() error {
	level, err := logrus.ParseLevel(t.config.LogLevel)
	if err != nil {
		return err
	}

	t.logger.SetLevel(level)
	return nil
}

func (t *TelegramBot) runBot() {
	// Get the TelegramToken environment variable
	telegramToken, exists := os.LookupEnv("TELEGRAM_TOKEN")
	if exists {
		t.logger.Info("token:" + telegramToken)
	}
	offset := 0
	//botURL := t.config.Telegram.APIURL + telegramToken

	// Start of new bot implementation
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		t.logger.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(offset)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	fmt.Print(".")
	for update := range updates {
		if update.CallbackQuery != nil {
			fmt.Print(update)

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Done"))
			writeToSerial(t, "open")
		}
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ну прівєтікі")
			switch update.Message.Text {
			case "help", "?", "/help", "/start":
				msg.ReplyMarkup = numericKeyboard
			case "open", "/open":
				msg.Text = "Ага блять, щас уже біжу..."
			}

			bot.Send(msg)
		}
	}

	// for {
	// 	updates, err := getUpdates(botURL, offset)
	// 	if err != nil {
	// 		t.logger.Error("Smth went wrong")
	// 	}
	// 	for _, update := range updates {

	// 		t.logger.Info(update)
	// 		//t.logger.Info(update.Message.Text)

	// 		//t.logger.Info(update.Message.From.Username)

	// 		err = respond(botURL, update)
	// 		offset = update.UpdateID + 1
	// 	}
	// }

	//writeToSerial(t, "open")

}

func getUpdates(botURL string, offset int) ([]Update, error) {
	resp, err := http.Get(botURL + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func respond(botURL string, update Update) error {
	var botMessage BotMessage
	botMessage.ChatID = update.Message.Chat.ChatID
	botMessage.Text = update.Message.From.Username + " хуйца сосни"

	h := json.RawMessage(`{
		"chat_id": 235023663,
		"text": "LOL"
		}`)

	c := struct {
		*json.RawMessage
	}{&h}

	b, err := json.MarshalIndent(&c, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}

	// buf, err := json.Marshal(botMessage)

	// res2D := &response2{
	//     Page:   1,
	//     Fruits: []string{"apple", "peach", "pear"}}

	// res2B, _ := json.Marshal(res2D)
	fmt.Println(string(b))

	if err != nil {
		return err
	}
	resp, err := http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(b))
	fmt.Println(resp)

	if err != nil {
		return err
	}
	return nil
}

func writeToSerial(t *TelegramBot, s string) {
	c := &serial.Config{Name: t.config.Serial.SerialName, Baud: t.config.Serial.SerialPort}
	sr, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	_, err = sr.Write([]byte(s))
	if err != nil {
		log.Fatal(err)
	}

}
