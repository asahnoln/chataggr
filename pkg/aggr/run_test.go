package aggr_test

import (
	"testing"
	"time"

	"github.com/asahnoln/chataggr/pkg/aggr"
)

func TestRun(t *testing.T) {
	c := make(chan aggr.Message)
	go aggr.Run([]aggr.Receiver{
		&stubSlowMessager{},
		&stubMessager{msgs: []aggr.Message{{}, {}}},
		&stubMessager{msgs: []aggr.Message{{}}},
	}, c)

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

	if got, want := len(msgs), 3; got != want {
		t.Errorf("want msg len %v, got %v", want, got)
	}
}

type stubMessager struct {
	msgs []aggr.Message
}

func (m *stubMessager) Receive(c chan aggr.Message) {
	for _, msg := range m.msgs {
		c <- msg
	}
}

type stubSlowMessager struct{}

func (m *stubSlowMessager) Receive(c chan aggr.Message) {
	time.Sleep(1 * time.Second)
}
