package receivers

import (
	"log"

	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/gempir/go-twitch-irc/v4"
)

type Twitch struct {
	clt *twitch.Client
}

func NewTwitch(clt *twitch.Client) *Twitch {
	return &Twitch{clt}
}

func (r *Twitch) Receive(c chan aggr.Message) {
	// r.clt.Join("something")
	err := r.clt.Connect()
	if err != nil {
		log.Fatalf("twitch client connect err: %v", err)
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
