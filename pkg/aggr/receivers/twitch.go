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
	startMsgs := []string{
		"CAP REQ :twitch.tv/tags twitch.tv/Command",
		"PASS SCHMOOPIIE",
		"NICK justinfan8865",
		"USER justinfan8865 8 * :justinfan8865",
	}
	for _, msg := range startMsgs {
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}

	c <- aggr.Message{
		Text: "Hi chat",
		User: "TwitchDev",
	}
	c <- aggr.Message{
		Text: "Hi chat",
		User: "TwitchDev",
	}
}
