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
	// FIX: gotiktoklive doesn't work anymore
	// rtt := createTikTokReceiver()
	simr := createStubInMemoryReceiver()

	c := make(chan aggr.Message)
	aggr.Run([]aggr.Receiver{simr}, c)

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

func createStubInMemoryReceiver() *stubInMemoryReceiver {
	return &stubInMemoryReceiver{}
}
