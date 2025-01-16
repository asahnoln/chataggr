package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

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
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	// simr := createStubInMemoryReceiver()

	// FIX: gotiktoklive doesn't work anymore
	// rtt := createTikTokReceiver()
	twr := createTwitchReceiver()

	c := make(chan aggr.Message)
	aggr.Run([]aggr.Receiver{twr}, c)

loop:
	for {
		select {
		case m := <-c:
			fmt.Printf("%s [%s]: %s\n", m.User, time.Now().Format(time.TimeOnly), m.Text)
		case <-interrupt:
			break loop
		}
	}
}

// FIX: gotiktoklive doesn't work anymore
func createTikTokReceiver() *receivers.TikTok {
	return receivers.NewTikTok("tiktok url")
}

func createTwitchReceiver() *receivers.Twitch {
	return receivers.NewTwitch("wss://irc-ws.chat.twitch.tv/")
}

func createStubInMemoryReceiver() *stubInMemoryReceiver {
	return &stubInMemoryReceiver{}
}
