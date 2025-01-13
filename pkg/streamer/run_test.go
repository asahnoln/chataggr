package streamer_test

import (
	"testing"
	"time"

	"github.com/asahnoln/streamer-chat-go/pkg/streamer"
)

func TestRun(t *testing.T) {
	c := make(chan streamer.Message)
	go streamer.Run([]streamer.Messager{
		&stubSlowMessager{},
		&stubMessager{msgs: []streamer.Message{{}, {}}},
		&stubMessager{msgs: []streamer.Message{{}}},
	}, c)

	msgs := []streamer.Message{}
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
	msgs []streamer.Message
}

func (m *stubMessager) Receive(c chan streamer.Message) {
	for _, msg := range m.msgs {
		c <- msg
	}
}

type stubSlowMessager struct{}

func (m *stubSlowMessager) Receive(c chan streamer.Message) {
	time.Sleep(1 * time.Second)
}
