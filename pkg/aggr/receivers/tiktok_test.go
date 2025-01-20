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
	r := receivers.NewTikTok("https://webcast.tiktok.com/webcast/im/fetch/?version_code=180800&device_platform=web&cookie_enabled=true&screen_width=1920&screen_height=1200&browser_language=en-US&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0%20(Windows)&browser_online=true&tz_name=Asia/Qyzylorda&aid=1988&app_name=tiktok_web&live_id=12&version_code=270000&debug=false&app_language=en&room_id=7461984576673843974&identity=audience&history_comment_count=6&fetch_rule=1&last_rtt=-1&cursor=0&internal_ext=0&sup_ws_ds_opt=1&resp_content_type=protobuf&did_rule=3&msToken=awJxSgYETZELnLx0W_4t_n1Ku6AqI3pr7uOwzG0nb5HFPeKinZ1Bo_1N72nDdVDlJ-8ARb7yQum-Jo4fMTTex3yzgGhJxwi78PTu7mdZk5p1WnYdhW6moImR4ZsVIoxaFnipmZcMENDRIJEmgIKV7A==&X-Bogus=DFSzsIVO42JANa7UtpO2Gw2J46s3&_signature=_02B4Z6wo00001IdK92wAAIDDaP0wlbVNcxCHSPPAAEZh2f")
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
	r := receivers.NewTikTok("https://webcast.tiktok.com/webcast/im/fetch/?version_code=180800&device_platform=web&cookie_enabled=true&screen_width=1920&screen_height=1200&browser_language=en-US&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0%20(Windows)&browser_online=true&tz_name=Asia/Qyzylorda&aid=1988&app_name=tiktok_web&live_id=12&version_code=270000&debug=false&app_language=en&room_id=7461277983610473221&identity=audience&history_comment_count=6&fetch_rule=1&last_rtt=-1&cursor=0&internal_ext=0&sup_ws_ds_opt=1&resp_content_type=protobuf&did_rule=3&msToken=eNoZdun-hYEZBmGb0j0PVuB2JelKkytdjE7yMowNdUp6HzHQuIbeebJcD1KmOE-uAaqJI5pvHPOpylMNERkkIOv4qnCerZ-xBtjLQlaSIkr3R4xVXS3tFnJvf4t9e6nIsSiSQWFIynEPnDk44jobzw==&X-Bogus=DFSzsIVOIGzANa7UtpSl7w2J46BG&_signature=_02B4Z6wo00001JKXL9AAAIDDfSDoKESQ7fCSlStAAEMy35")

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
