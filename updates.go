package main

import (
	"fmt"
	"time"
	"os"
	"net/http"
	"strings"
	"net/url"
	"strconv"
)

var activeChats = make(map[int]int)
var userQ = make(chan int, 200)

type botMessage struct {
	to   int
	text string
}

var messageQ = make(chan *botMessage, 15)

func processUpdate(update updateObject) {
	switch update.Type {
	case "message_new":
		userId := update.Object.UserID
		message := update.Object.Body
		fmt.Println(userId, message)
		if _, inChat := activeChats[userId]; inChat {
			if message == "!!!" {
				sendVkMessage(userId, "⚡ Беседа окончена")
				sendVkMessage(activeChats[userId], "⚡ Беседа окончена")
				removeFromActiveChat(userId)
				return
			}

			sendVkMessage(activeChats[userId], "[собеседник️]: "+message)
			return
		}

		switch message {
		case "го":
			userQ <- userId
			sendVkMessage(userId, "⚡ Ты в очереди...")
			return

		default:
			sendVkMessage(userId, `
⚡ Чтобы найти собеседника отправь "го"
⚡ Чтобы завершить чат с собеседником оправь в чат "!!!"
`)
		}

	}
}
func startChatLoop() {
	go func() { // creating chats from queued users
		for {
			time.Sleep(time.Millisecond * 10)
			u1, u2 := <-userQ, <-userQ
			createChat(u1, u2)
		}
	}()
	go func() { // send queued messages
		for {
			time.Sleep(time.Millisecond * 100)
			message := <-messageQ
			postMessageActual(message)
		}
	}()
}

func removeFromActiveChat(u1 int) {
	if partner, inChat := activeChats[u1]; inChat {
		delete(activeChats, u1)
		delete(activeChats, partner)
	}
}

func createChat(user1, user2 int) {
	removeFromActiveChat(user1)
	removeFromActiveChat(user2)
	activeChats[user1] = user2
	activeChats[user2] = user1
	sendVkMessage(user1, "⚡ Собеседник найден, общайтесь!")
	sendVkMessage(user2, "⚡ Собеседник найден, общайтесь!")
}

func sendVkMessage(to int, text string) {
	message := &botMessage{to, text}
	messageQ <- message
}
func postMessageActual(message *botMessage) {
	//var postUrl = fmt.Sprintf("https://api.vk.com/method/messages.send?user_id=%d&access_token=%s&v=%s&message=%s", to, os.Getenv("TOKEN"), "5.74", text)
	data := url.Values{
		"user_id":      {strconv.Itoa(message.to)},
		"access_token": {os.Getenv("TOKEN")},
		"v":            {"5.74"},
		"message":      {message.text}}

	r, _ := http.DefaultClient.Post("https://api.vk.com/method/messages.send", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	var body []byte
	r.Body.Read(body)
	fmt.Println(string(body))
}
