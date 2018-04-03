package main

import (
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"errors"
	"os"
)
var getLongpollServerAPI = fmt.Sprintf("https://api.vk.com/method/groups.getLongPollServer?group_id=%d&access_token=%s&v=%s", groupID, os.Getenv("TOKEN"), "5.74")

func longpoll(server, key, ts string, updatesChan chan updateObject) {
	for {
		log.Print("poll...")
		r, _ := http.DefaultClient.Get(fmt.Sprintf("%s?act=a_check&key=%s&ts=%s&wait=10",
			server, key, ts))
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print("cannot read response body")
			close(updatesChan)
			break
		}

		var update longpollUpdate
		if e := json.Unmarshal(b, &update); e == nil {
			ts = update.Ts
			for _, u := range update.Updates {
				updatesChan <- u
			}
			continue
		}

		var lpError failedLP
		json.Unmarshal(b, &lpError)
		fmt.Println(string(b))
		if lpError.Failed == 1 {
			ts = strconv.Itoa(lpError.Ts)
			continue
		} else {
			server, key, ts, err = getLongpollData()
			if err != nil {
				log.Print("polling failed")
				close(updatesChan)
				break
			}
			continue
		}
	}
}
func longPollUpdates() (<-chan updateObject) {
	server, key, ts, err := getLongpollData()
	if err != nil {
		log.Fatal("cannot longpoll :(", err)
	}

	updates := make(chan updateObject)
	go longpoll(server, key, ts, updates)
	return updates
}

func getLongpollData() (string, string, string, error) {
	r, e := http.DefaultClient.Get(getLongpollServerAPI)
	if e != nil {
		return "", "", "", e
	}

	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return "", "", "", e
	}

	var response longpollResponse
	json.Unmarshal(b, &response)
	if response.Error == nil && response.Response != nil {
		return response.Response.Server, response.Response.Key, strconv.Itoa(response.Response.Ts), nil
	}

	return "", "", "", errors.New(fmt.Sprintf("error during getting server: [%d] %s",
		response.Error.ErrorCode,
		response.Error.ErrorMsg))
}

