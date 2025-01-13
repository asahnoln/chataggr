package receivers_test

import (
	"testing"
	"time"

	"github.com/Davincible/gotiktoklive"
	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/asahnoln/chataggr/pkg/aggr/receivers"
)

// FIX: gotiktoklive doesn't work anymore
func TestTiktok(t *testing.T) {
	l := &gotiktoklive.Live{Events: make(chan interface{})}
	go func() {
		l.Events <- gotiktoklive.ChatEvent{
			Comment: "Hi Tik",
			User:    &gotiktoklive.User{Nickname: "someone"},
		}
		l.Events <- gotiktoklive.ChatEvent{Comment: "Hi Tik 2"}
	}()

	c := make(chan aggr.Message)
	r := receivers.NewTikTok(l)

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
		t.Errorf("want %v, got %v", want, got)
	}

	if got, want := msgs[0].Text, "Hi Tik"; got != want {
		t.Errorf("want %v, got %v", want, got)
	}

	if got, want := msgs[0].User, "someone"; got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
