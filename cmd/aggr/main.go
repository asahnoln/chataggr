package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Davincible/gotiktoklive"
	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/asahnoln/chataggr/pkg/aggr/receivers"
)

// stubInMemoryReceiver is used for debugging purposes
type stubInMemoryReceiver struct{}

func (r *stubInMemoryReceiver) Receive(c chan aggr.Message) {
	t := time.Tick(3 * time.Second)
	for range t {
		c <- aggr.Message{User: "asahnoln", Text: "Hi"}
	}
}

func main() {
	// simr := createStubInMemoryReceiver()

	// FIX: gotiktoklive doesn't work anymore
	// rtt := createTikTokReceiver()
	twr := createTwitchReceiver()

	c := make(chan aggr.Message)
	aggr.Run([]aggr.Receiver{twr}, c)

	for m := range c {
		fmt.Printf("%s [%s]: %s\n", m.User, time.Now().Format(time.TimeOnly), m.Text)
	}
}

// FIX: gotiktoklive doesn't work anymore
func createTikTokReceiver() *receivers.TikTok {
	tt := gotiktoklive.NewTikTok()
	live, err := tt.TrackUser("asahnoln")
	if err != nil {
		log.Panicf("track user error: %v", err)
	}

	return receivers.NewTikTok(live)
}

func createTwitchReceiver() *receivers.Twitch {
	return receivers.NewTwitch("wss://irc-ws.chat.twitch.tv/")
}

func createStubInMemoryReceiver() *stubInMemoryReceiver {
	return &stubInMemoryReceiver{}
}
