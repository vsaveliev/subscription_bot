package main

import (
	"encoding/json"
	"gopkg.in/telegram-bot-api.v3"
	"log"
	"net/http"
	"strconv"
)

var userSessions map[string]UserSession

type Message struct {
	Message string
}

func main() {
	bot, err := createBot()
	if err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/subscription/add", func(w http.ResponseWriter, r *http.Request) {
		var msg Message

		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		sendVacationToAll(bot, msg.Message)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	})

	go initApp(bot)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initApp(bot *tgbotapi.BotAPI) {
	userSessions = getUserSessions()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		processUpdate(bot, update)
	}
}

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	//	switch {
	//	case msg.Text == "/start":
	_, isExists := userSessions[strconv.FormatInt(update.Message.Chat.ID, 10)]

	if !isExists {
		sendGreetingMessage(bot, msg)
		addSession(update)
	}
	//	}
}

func createBot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI("258957984:AAGmMmfsA8eYeHx8OB_mvyFqHZPGoOxOvds")
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot, nil
}

// TODO: create special struct for it #1
func sendGreetingMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) (bool, error) {
	msg.Text = "Привет. Я твой персональный помощник по поиску вакансий! Теперь я буду отправлять вам вакансии"

	_, err := bot.Send(msg)
	if err != nil {
		return false, err
	}

	return true, nil
}

func sendVacationToAll(bot *tgbotapi.BotAPI, message string) {

	log.Printf("Send message to all - %s", message)

	for chatId, session := range userSessions {
		log.Printf("Send message to - %s", chatId)

		msg := tgbotapi.MessageConfig{
			Text: message,
			BaseChat: tgbotapi.BaseChat{
				ChatID: session.ChatId,
			},
			ParseMode: "HTML",
		}

		_, err := bot.Send(msg)
		if err != nil {
			return
		}
	}
}
