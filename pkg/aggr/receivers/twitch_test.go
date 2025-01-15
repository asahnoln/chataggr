package receivers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/asahnoln/chataggr/pkg/aggr/receivers"
	"github.com/gempir/go-twitch-irc/v4"
)

func TestTwitch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("req %v", r.URL)
		t.Log("AAAAAAAAAAAAAAAAAAA")
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Logf("error reading request body: %v", err)
		}

		defer r.Body.Close()

		t.Logf("webserver %s", b)
	}))
	defer srv.Close()

	clt := twitch.NewAnonymousClient()
	clt.IrcAddress = strings.TrimPrefix(srv.URL, "http://")
	clt.TLS = false

	c := make(chan aggr.Message)
	r := receivers.NewTwitch(clt)

	aggr.Run([]aggr.Receiver{r}, c)

	msgs := []aggr.Message{}
	timer := time.NewTimer(5 * time.Second)

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
		t.Fatalf("want len %v, got %v", want, got)
	}

	if got, want := msgs[0].Text, "Hi chat"; got != want {
		t.Errorf("want text %v, got %v", want, got)
	}

	if got, want := msgs[0].User, "TwitchDev"; got != want {
		t.Errorf("want name %v, got %v", want, got)
	}
}
