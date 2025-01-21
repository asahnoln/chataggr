package receivers

import (
	"log"
	"strings"

	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/gorilla/websocket"
)

var twitchConnectMsgs [5]string = [5]string{
	"CAP REQ :twitch.tv/tags twitch.tv/Command",
	"PASS SCHMOOPIIE",
	"NICK justinfan8865",
	"USER justinfan8865 8 * :justinfan8865",
	"JOIN #asahnoln",
}

type Twitch struct {
	url string
}

func NewTwitch(url string) *Twitch {
	return &Twitch{url}
}

func (r *Twitch) Receive(c chan aggr.Message) {
	conn := connectToTwitch(r.url)
	defer conn.Close()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ws read err: %v", err)
			return
		}

		msg := string(data)
		text := findSubstrBetween(msg, "PRIVMSG #asahnoln :", "\r\n")
		if text == "" {
			continue
		}

		c <- aggr.Message{
			Receiver: r,
			Text:     text,
			User:     findSubstrBetween(msg, "display-name=", ";"),
		}
	}
}

func connectToTwitch(url string) *websocket.Conn {
	// FIX: Handle error
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("ws dial err: %v", err)
	}

	sendConnectMsgsToTwitch(conn)
	return conn
}

func sendConnectMsgsToTwitch(conn *websocket.Conn) {
	for _, msg := range twitchConnectMsgs {
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func findSubstrBetween(s, fromSubstr, toSubstr string) string {
	start := strings.Index(s, fromSubstr)
	if start == -1 {
		return ""
	}

	start += len(fromSubstr)
	end := strings.Index(s[start:], toSubstr)
	return s[start:][:end]
}
