package receivers_test

import (
	"testing"
	"time"

	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/asahnoln/chataggr/pkg/aggr/receivers"
)

func TestTwitch(t *testing.T) {
	c := make(chan aggr.Message)
	r := receivers.NewTwitch(&stubTwitchWSConn{})

	aggr.Run([]aggr.Receiver{r}, c)

	msgs := []aggr.Message{}
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
		t.Errorf("want len %v, got %v", want, got)
	}

	if got, want := msgs[0].Text, "Hi chat"; got != want {
		t.Errorf("want text %v, got %v", want, got)
	}

	if got, want := msgs[0].User, "TwitchDev"; got != want {
		t.Errorf("want name %v, got %v", want, got)
	}
}

type stubTwitchWSConn struct{}

func (r *stubTwitchWSConn) Connect(c chan []byte) {
	c <- []byte(`{
      "payload": {
        "subscription": {
          "type": "channel.chat.message"
        },
        "event": {
          "chatter_user_name": "TwitchDev",
          "message": {
            "text": "Hi chat"
          }
        }
      }
  }`)
	c <- []byte(`{
      "payload": {
        "subscription": {
          "type": "bla.bla"
        },
        "event": {
          "blabla": 1
        }
      }
  }`)
	c <- []byte(`{
      "payload": {
        "subscription": {
          "type": "channel.chat.message"
        },
        "event": {
          "chatter_user_name": "TwitchDev",
          "message": {
            "text": "Hi chat"
          }
        }
      }
  }`)
}
