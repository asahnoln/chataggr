package receivers

import (
	"log"
	"strings"

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

	startMsgs := []string{
		"CAP REQ :twitch.tv/tags twitch.tv/Command",
		"PASS SCHMOOPIIE",
		"NICK justinfan8865",
		"USER justinfan8865 8 * :justinfan8865",
		"JOIN #asahnoln",
	}
	for _, msg := range startMsgs {
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ws read err: %v", err)
			return
		}

		msg := string(data)
		c <- aggr.Message{
			Text: findValue(msg, "PRIVMSG #asahnoln :", ""),
			User: findValue(msg, "display-name=", ";"),
		}
	}
}

func findValue(msg string, name string, sep string) string {
	start := strings.Index(msg, name) + len(name)

	if sep == "" {
		return msg[start:]
	}

	end := strings.Index(msg[start:], sep)
	return msg[start:][:end]
}
