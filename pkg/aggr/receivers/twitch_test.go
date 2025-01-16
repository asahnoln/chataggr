package receivers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/asahnoln/chataggr/pkg/aggr/receivers"
	"github.com/gorilla/websocket"
)

func TestTwitchWS(t *testing.T) {
	upg := websocket.Upgrader{}
	srvCalled := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testMsgs := []string{
			"CAP REQ :twitch.tv/tags twitch.tv/Command",
			"PASS SCHMOOPIIE",
			"NICK justinfan8865",
			"USER justinfan8865 8 * :justinfan8865",
			"JOIN #asahnoln",
		}
		srvCalled = true
		c, _ := upg.Upgrade(w, r, nil)
		defer c.Close()

		// TODO: for loop for incoming messages
		for _, tm := range testMsgs {
			mt, message, _ := c.ReadMessage()
			if got, want := mt, websocket.TextMessage; got != want {
				t.Errorf("srv want msg %v, got %v", want, got)
			}
			if got, want := string(message), tm; got != want {
				t.Errorf("srv want msg %v, got %v", want, got)
			}
		}
	}))
	defer srv.Close()

	tw := receivers.NewTwitch(strings.Replace(srv.URL, "http", "ws", 1))
	c := make(chan aggr.Message, 100)
	tw.Receive(c)

	if got, want := srvCalled, true; got != want {
		t.Fatal("want to call the server, but it didn't receive a request")
	}
}

func TestTwitchMessages(t *testing.T) {
	upg := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := upg.Upgrade(w, r, nil)
		defer c.Close()

		c.WriteMessage(websocket.TextMessage, []byte("@badge-info=;badges=broadcaster/1;client-nonce=30620c76bc7d8ba00c9fe4c4d81c5ef3;color=;display-name=TwitchDev;emotes=;first-msg=0;flags=;id=6fdb1223-9252-4ee3-b4dc-c0ac88c2b372;mod=0;returning-chatter=0;room-id=39182089;subscriber=0;tmi-sent-ts=1736906306256;turbo=0;user-id=39182089;user-type= :asahnoln!asahnoln@asahnoln.tmi.twitch.tv PRIVMSG #asahnoln :Hi chat"))
		c.WriteMessage(websocket.TextMessage, []byte("@badge-info=;badges=broadcaster/1;client-nonce=30620c76bc7d8ba00c9fe4c4d81c5ef3;color=;display-name=Asahnoln;emotes=;first-msg=0;flags=;id=6fdb1223-9252-4ee3-b4dc-c0ac88c2b372;mod=0;returning-chatter=0;room-id=39182089;subscriber=0;tmi-sent-ts=1736906306256;turbo=0;user-id=39182089;user-type= :asahnoln!asahnoln@asahnoln.tmi.twitch.tv PRIVMSG #asahnoln :Hey there!"))
	}))
	defer srv.Close()

	c := make(chan aggr.Message)
	r := receivers.NewTwitch(strings.Replace(srv.URL, "http", "ws", 1))

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
		t.Fatalf("want len %v, got %v", want, got)
	}

	if got, want := msgs[0].Text, "Hi chat"; got != want {
		t.Errorf("want text %v, got %v", want, got)
	}

	if got, want := msgs[0].User, "TwitchDev"; got != want {
		t.Errorf("want name %v, got %v", want, got)
	}
}
