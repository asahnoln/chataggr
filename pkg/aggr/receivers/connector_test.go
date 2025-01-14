package receivers_test

import (
	"testing"
	"time"

	"github.com/asahnoln/chataggr/pkg/aggr/receivers"
)

// Connect to ws
// wss://eventsub.wss.twitch.tv/ws
func TestTwitchConnector(t *testing.T) {
	var conn receivers.Connector = receivers.NewTwitchConnector()

	c := make(chan []byte)
	conn.Connect(c)

	msgs := [][]byte{}
	timer := time.NewTimer(1 * time.Millisecond)

l:
	for {
		select {
		case m := <-c:
			msgs = append(msgs, m)
		case <-timer.C:
			break l
		}
	}

	if got, want := len(msgs), 2; got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
