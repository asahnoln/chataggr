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
	r := receivers.NewTikTok("https://webcast.tiktok.com/webcast/im/fetch/?version_code=180800&device_platform=web&cookie_enabled=true&screen_width=1920&screen_height=1200&browser_language=en-US&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0%20(Windows)&browser_online=true&tz_name=Asia/Qyzylorda&aid=1988&app_name=tiktok_web&live_id=12&version_code=270000&debug=false&app_language=en&room_id=7462401799221185285&identity=audience&history_comment_count=6&fetch_rule=2&last_rtt=965&internal_ext=fetch_time:1737476887208%7Cstart_time:0%7Cack_ids:,,%7Cflag:0%7Cseq:1%7Cnext_cursor:1737476887208_7462406085098340604_1_7462406390041018425_7462403813060654792_0%7Cwss_info:0-1737476887208-0-0&cursor=1737476887208_7462406085098340604_1_7462406390041018425_7462403813060654792_0&history_comment_cursor=7462402946701396742&sup_ws_ds_opt=1&resp_content_type=protobuf&did_rule=3&msToken=ZSf430xg7a8o0NoMy9YprkwshQJ3c304zV0YD73uEq9a0A3fXwJsXU-mcEVCJzQ28tADMon83l5ablPbSCPazPvAEvDmklBqmYlv8hJtxxpXCl0KqTFkjrgKLky5CzausNEPOELTVQ4d&X-Bogus=DFSzsIVYSLXANcYTtpuMAd2J46BM&_signature=_02B4Z6wo00001oGmDxQAAIDBbhHI7sqNMz6BpAuAAMfy4a")
	r.Receive(c)
}

func TestTiktok(t *testing.T) {
	upg := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := upg.Upgrade(w, r, nil)
		defer c.Close()

		c.WriteMessage(websocket.TextMessage, []byte("Hi Tik,someone"))
		c.WriteMessage(websocket.TextMessage, []byte("Hi Tok,anyone"))
	}))
	defer srv.Close()

	c := make(chan aggr.Message)
	r := receivers.NewTikTok("https://webcast.tiktok.com/webcast/im/fetch/?version_code=180800&device_platform=web&cookie_enabled=true&screen_width=1920&screen_height=1200&browser_language=en-US&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0%20(Windows)&browser_online=true&tz_name=Asia/Qyzylorda&aid=1988&app_name=tiktok_web&live_id=12&version_code=270000&debug=false&app_language=en&room_id=7462401799221185285&identity=audience&history_comment_count=6&fetch_rule=1&last_rtt=-1&cursor=0&internal_ext=0&sup_ws_ds_opt=1&resp_content_type=protobuf&did_rule=3&msToken=ZSf430xg7a8o0NoMy9YprkwshQJ3c304zV0YD73uEq9a0A3fXwJsXU-mcEVCJzQ28tADMon83l5ablPbSCPazPvAEvDmklBqmYlv8hJtxxpXCl0KqTFkjrgKLky5CzausNEPOELTVQ4d&X-Bogus=DFSzsIVYJyGANcYTtpuMAw2J46Bj&_signature=_02B4Z6wo00001bO-eRwAAIDCXAm-5EKYQHGzvH2AAAth3a")

	aggr.Run([]aggr.Receiver{r}, c)

	msgs := []aggr.Message{}
	timer := time.NewTimer(10 * time.Second)

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
