package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/telegram-bot-api.v3"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type UserSession struct {
	ChatId int64         `json:"chat_id"`
	User   tgbotapi.User `json:"user"`
}

func getUserSessions() map[string]UserSession {
	raw, err := ioutil.ReadFile("/var/www/gocode/src/github.com/user/mypromoagent/sessions.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	c := make(map[string]UserSession)
	json.Unmarshal(raw, &c)
	return c
}

func saveUserSessions(sessions map[string]UserSession) {
	j, jerr := json.MarshalIndent(sessions, "", "  ")
	if jerr != nil {
		fmt.Println("jerr:", jerr.Error())
	}

	werr := ioutil.WriteFile("./sessions.json", j, os.ModeExclusive)
	if werr != nil {
		fmt.Println("werr:", werr.Error())
	}
}

func addSession(update tgbotapi.Update) {
	chatId := strconv.FormatInt(update.Message.Chat.ID, 10)

	newSession := UserSession{
		ChatId: update.Message.Chat.ID,
		User:   update.Message.From,
	}
	userSessions[chatId] = newSession

	go saveUserSessions(userSessions)

	log.Printf("new session =", newSession.ChatId)

}
