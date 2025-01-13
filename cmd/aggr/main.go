package main

import (
	"fmt"
	"log"

	"github.com/Davincible/gotiktoklive"
	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/asahnoln/chataggr/pkg/aggr/receivers"
)

func main() {
	tt := gotiktoklive.NewTikTok()
	live, err := tt.TrackUser("asahnoln")
	if err != nil {
		log.Panicf("track user error: %v", err)
	}

	rtt := receivers.NewTikTok(live)

	c := make(chan aggr.Message)
	aggr.Run([]aggr.Receiver{rtt}, c)

	for m := range c {
		fmt.Printf("%s: %s\n", m.User, m.Text)
	}
}
