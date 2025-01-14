package receivers

import (
	"encoding/json"
	"log"

	"github.com/asahnoln/chataggr/pkg/aggr"
)

type Twitch struct {
	conn Connector
}

type TwitchWSNotification struct {
	Payload struct {
		Subscription struct {
			Type string `json:"type"`
		} `json:"subscription"`
		Event struct {
			ChatterUserName string `json:"chatter_user_name"`
			Message         struct {
				Text string `json:"text"`
			} `json:"message"`
		} `json:"event"`
	} `json:"payload"`
}

func NewTwitch(conn Connector) *Twitch {
	return &Twitch{conn}
}

func (r *Twitch) Receive(c chan aggr.Message) {
	bodyChan := make(chan []byte)
	go r.conn.Connect(bodyChan)

	for b := range bodyChan {
		var resp TwitchWSNotification

		// FIX: Handle error
		err := json.Unmarshal(b, &resp)
		if err != nil {
			log.Printf("twitch unmarshal err: %v", err)
		}

		if resp.Payload.Subscription.Type != "channel.chat.message" {
			continue
		}

		c <- aggr.Message{
			Text: resp.Payload.Event.Message.Text,
			User: resp.Payload.Event.ChatterUserName,
		}
	}
}
