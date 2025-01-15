package receivers

import (
	"log"

	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/gorilla/websocket"
)

type Twitch struct {
	url string
}

func NewTwitch(url string) *Twitch {
	return &Twitch{url}
}

func (r *Twitch) Receive(c chan aggr.Message) {
	// FIX: Handle error
	conn, _, err := websocket.DefaultDialer.Dial(r.url, nil)
	if err != nil {
		log.Fatalf("ws dial err: %v", err)
	}
	defer conn.Close()

	// TODO: go func reading incoming messages
	conn.WriteMessage(websocket.TextMessage, []byte("PASS SCHMOOPIIE"))

	c <- aggr.Message{
		Text: "Hi chat",
		User: "TwitchDev",
	}
	c <- aggr.Message{
		Text: "Hi chat",
		User: "TwitchDev",
	}
}
