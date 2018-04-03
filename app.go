package main

import (
	"log"
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

var groupID = 71132642

func main() {
	go startChatLoop()
	//for u := range longPollUpdates() {
	//	go processUpdate(u)
	//}

	http.HandleFunc("/", postHandler)
	log.Println("Listening...")
	http.ListenAndServe(":5555", nil)

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	var update updateObject
	json.Unmarshal(bodyBytes, &update)
	if update.Secret != "sdlkfhk89394"{
		return
	}
	go processUpdate(update)
	fmt.Fprintf(w, "ok")
}


