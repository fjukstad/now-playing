package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var playing string

func live(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		artist, track := GetCurrent()
		nowPlaying := artist + "," + track
		if playing != nowPlaying {
			err = c.WriteMessage(1, []byte(nowPlaying))
			if err != nil {
				log.Println("write:", err)
				break
			}
			playing = nowPlaying
		}
		time.Sleep(1 * time.Second)
	}
}
