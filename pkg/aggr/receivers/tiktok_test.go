package receivers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/asahnoln/chataggr/pkg/aggr/receivers"
	"github.com/gorilla/websocket"
)

func TestTikTokIntegration(t *testing.T) {
	c := make(chan aggr.Message, 100)
	r := receivers.NewTikTok("wss://webcast16-ws-useast1a.tiktok.com/webcast/im/ws_proxy/ws_reuse_supplement/?version_code=180800&device_platform=web&cookie_enabled=true&screen_width=1920&screen_height=1200&browser_language=en-US&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0%20(Windows)&browser_online=true&tz_name=Asia/Qyzylorda&app_name=tiktok_web&sup_ws_ds_opt=1&version_code=270000&update_version_code=2.0.0&compress=gzip&wrss=zJWkjhmSlSzaZ_mWfao7806LeYcdYxfYtrOnWr1gRL0&host=https://webcast.tiktok.com&aid=1988&live_id=12&debug=false&app_language=en&room_id=7461277983610473221&identity=audience&history_comment_count=6&heartbeat_duration=0&last_rtt=641&internal_ext=fetch_time:1737215424271|start_time:0|ack_ids:,,|flag:0|seq:1|next_cursor:1737215424271_7461283264683049335_1_7461282775056777256_7461279944673332788_0|wss_info:0-1737215424271-0-0&cursor=1737215424271_7461283264683049335_1_7461282775056777256_7461279944673332788_0&history_comment_cursor=7461278899575687942&resp_content_type=protobuf&did_rule=3")
	r.Receive(c)
}

func TestTiktok(t *testing.T) {
	t.SkipNow()
	upg := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := upg.Upgrade(w, r, nil)
		defer c.Close()

		c.WriteMessage(websocket.TextMessage, []byte("Hi Tik,someone"))
		c.WriteMessage(websocket.TextMessage, []byte("Hi Tok,anyone"))
	}))
	defer srv.Close()

	c := make(chan aggr.Message)
	r := receivers.NewTikTok("wss://webcast16-ws-useast1a.tiktok.com/webcast/im/ws_proxy/ws_reuse_supplement/?version_code=180800&device_platform=web&cookie_enabled=true&screen_width=1920&screen_height=1200&browser_language=en-US&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0%20(Windows)&browser_online=true&tz_name=Asia/Qyzylorda&app_name=tiktok_web&sup_ws_ds_opt=1&version_code=270000&update_version_code=2.0.0&compress=gzip&wrss=zJWkjhmSlSzaZ_mWfao7806LeYcdYxfYtrOnWr1gRL0&host=https://webcast.tiktok.com&aid=1988&live_id=12&debug=false&app_language=en&room_id=7461277983610473221&identity=audience&history_comment_count=6&heartbeat_duration=0&last_rtt=641&internal_ext=fetch_time:1737215424271|start_time:0|ack_ids:,,|flag:0|seq:1|next_cursor:1737215424271_7461283264683049335_1_7461282775056777256_7461279944673332788_0|wss_info:0-1737215424271-0-0&cursor=1737215424271_7461283264683049335_1_7461282775056777256_7461279944673332788_0&history_comment_cursor=7461278899575687942&resp_content_type=protobuf&did_rule=3")

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
		t.Fatalf("want %v, got %v", want, got)
	}

	if got, want := msgs[0].Text, "Hi Tik"; got != want {
		t.Errorf("want %v, got %v", want, got)
	}

	if got, want := msgs[0].User, "someone"; got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
