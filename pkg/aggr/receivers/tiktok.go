package receivers

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/asahnoln/chataggr/pkg/aggr"
	chataggr "github.com/asahnoln/chataggr/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type TikTok struct {
	url string
}

func NewTikTok(url string) *TikTok {
	return &TikTok{url}
}

// Seems like I have to hardcode in first url and cookie
func (r *TikTok) Receive(c chan aggr.Message) {
	// TODO: fetch room id
	// TODO: fetch msToken
	// TODO: fetch X-Bogus
	// TODO: make signature
	mainURL, err := url.Parse(r.url)
	if err != nil {
		log.Fatalf("parsing webcast url err: %v", err)
	}

	// Create request for first main URL
	req, err := http.NewRequest(http.MethodGet, mainURL.String(), nil)
	if err != nil {
		log.Fatalf("webcast req err: %v", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:134.0) Gecko/20100101 Firefox/134.0")
	// TODO: fetch cookie?
	// WARN: IMPORTANT TO SEND COOKIE
	cookie := "ttwid=1%7C5XfWiTRds5mhaylBKPImrkN_-DXgtkG-B9RRMVMYXlE%7C1735837246%7C3c8891d53acf7a5b4ad1071797338b13c30f14a4b3c6a6b20212539b912bd2a4; tt_chain_token=80E7G904W83fcHsjnaWXiA==; msToken=tfQyFmr9GjVszUYX2c2hUsjDDa2gAv1E9pqcOUNPkUTrPrqfI-RmhxieMRHPWLgxv5SbKPcOlE80zHRudm_aYk7A_MoNWyxwJVtxiWaGt5coWUehB_jc9SPmdEZl0bhjEN557FvumlHe3Vs13pM7nA==; odin_tt=e926706ac66c5a2a62302f898508391000ad447a847627c7ee9262c0dc1123e16fde4a5530febd6810473e25916fd7815319b305ca339ca9b3e27c02c1fd520ffdbc05b8ca7b7a53c053bbf32b721708; passport_csrf_token=2d8e3e5d746bf7715f4a937e1da233c3; passport_csrf_token_default=2d8e3e5d746bf7715f4a937e1da233c3; multi_sids=6854322356408976389%3A39c7d85c0281664ce4fc7f1075f21568; cmpl_token=AgQQAPOsF-RO0o8hn3mF9F08_T51wNrW_4cnYNrqyQ; passport_auth_status=651f117f4ccc3d516fadff4b3d8e2715%2C; passport_auth_status_ss=651f117f4ccc3d516fadff4b3d8e2715%2C; sid_guard=39c7d85c0281664ce4fc7f1075f21568%7C1735837273%7C15552000%7CTue%2C+01-Jul-2025+17%3A01%3A13+GMT; uid_tt=eaef1ae9ffe38ee1a32b335e34a616019bdba1debe93085624e090fa9613ba63; uid_tt_ss=eaef1ae9ffe38ee1a32b335e34a616019bdba1debe93085624e090fa9613ba63; sid_tt=39c7d85c0281664ce4fc7f1075f21568; sessionid=39c7d85c0281664ce4fc7f1075f21568; sessionid_ss=39c7d85c0281664ce4fc7f1075f21568; sid_ucp_v1=1.0.0-KGEyODdhMzAxNjljNDA5Yjk0YzRjNzk2NTVkNDA4OGRiZmY1NTg3NzgKIQiFiJru4Lvcj18Q2YzbuwYYswsgDDD05P34BTgIQBJIBBADGgZtYWxpdmEiIDM5YzdkODVjMDI4MTY2NGNlNGZjN2YxMDc1ZjIxNTY4; ssid_ucp_v1=1.0.0-KGEyODdhMzAxNjljNDA5Yjk0YzRjNzk2NTVkNDA4OGRiZmY1NTg3NzgKIQiFiJru4Lvcj18Q2YzbuwYYswsgDDD05P34BTgIQBJIBBADGgZtYWxpdmEiIDM5YzdkODVjMDI4MTY2NGNlNGZjN2YxMDc1ZjIxNTY4; store-idc=maliva; store-country-code=kz; store-country-code-src=uid; tt-target-idc=useast1a; tt-target-idc-sign=hdPSGbtamqD5GvFwwuZEAoJvDzGVq5fbYOmVCRJvt69QkFjgdrIj_H5LE_h1hrs6Iz3bMyOOPseWkfEInbaSz2-gk0z6SEfvmFqdMAFYZUt2Knz-B5Bq50PnE_XeRs5pUqlll5pRKWANos5aTQTUz5EUCm4M9ZGSOi7AO3BNDtvN_h3kJhqtoy-C_OHSSMChWXELtmMpWm8Da8XHWVEV1trrayrPcF2ySN4jlkxr4m1jxUcervDFIgPUdmJ29TDU9H-aM74HRCBfR6jIgVKZ_3R7Dw9mYlhOQb7YE25A_YRZP4uMcxGWU05iYXUr8Z7hzwwUi3P9MD4f6dFnLYw_OqAHZHntAjvhic13GV6KFYu-nUWt8WufATGGAO7Y9D97TKUnEf8zKTOYTOrEitVIIUzbLup9EupF7QbGHlTaPoLuTeQyUuDdH6sGR4UodQtXm8LPtGvfhKcqovYtrGIsCPWfak9mOvS9bSK7MPDmAedDg-GWIe-9NwLeoq6QUNk7; tt_csrf_token=0vjm1tIz-p6M_LDzdB5Gmz1HEer5VfUXcKOQ; csrfToken=MKZEdlFH-3t3U6ouWGkID06FEe6_K1hQK228; s_v_web_id=verify_m5fkoyne_4Mce8XqP_d1uh_4w6l_BOTl_ucoxHn1KrPpt; csrf_session_id=319e32e5d82cbdfecd6c1e62ac4f67e3; ak_bmsc=C3B5C04FCEE054BD16C240BA8D68EB85~000000000000000000000000000000~YAAQjthUuHZ6cIqUAQAAYSiijRoYg3fZ7y9ZXtnlGHAFz9ysNM9hxdj28zoPHopVnTVIvmaYxtyuR+gPO+yG+rSJkfs7diT31vaO+zGpmdvt1meAHTkrmYIMhDVCTmGoixMpKvmbg6ct3RUXphnx9yLyQmnNO+tCqwl9qA5vmU1e18neZI5rD5xpK2O4yk/WRSPgBnRD9XQpNXHU7SH2p/hnaaAuoENhcDz0yDIZcnpW1tMrgnY619joFlUzuau3zpUPWEeI/SZQnGVl1VD3lSzZq/5U7MFKIH7hCn5CKpFQOGLV9QJ/ns64kDn9ttceUd8NsIh2NTkyCGyw6ofWqJywbbp9szqUiY2aYbtHW67IlPSPFagVFOHbx0JH9ky7W20ZNjLYQm4=; bm_sv=B98E7496E4508EB3B652162421ADFE80~YAAQVHkQAgcK0YqUAQAAhp/njRoR0P07/nlhGlak+FWRb+0GORm9Vn6/kPlqmZdYOdklzQA7sb1nAqdons3c+ZhFGwZ0Da0dOgi4gsS9RscBg+qKNsN3133fdEYqWkb3gaQzXzy0ymaI1jHuoOszOd/D5lAgI6Mgqy5HId/9lJhyVD614a73i919P3/YiqeGqcLqPuuBWnIa8894oVVYW9AEJhT8AD/bVUR0nfpnR7RqTNTpgLbOKrXfRfDPcccK~1"
	req.Header.Add("Cookie", cookie)

	// Make request to first main URL
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("webcast do err: %v", err)
	}

	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatalf("webcast body read err: %v", err)
	}

	var wb chataggr.WebcastResponse
	if err := proto.Unmarshal(b, &wb); err != nil {
		log.Fatalf("unmarshal proto err: %v", err)
	}

	// Prepare data for websocket request (URL is not real as it is gonna be replace with WsURL from webcast response)
	// WARN: For now taking URL straight from browser
	u, err := url.Parse("wss://replaceWithWrssParam.com/?version_code=180800&device_platform=web&cookie_enabled=true&screen_width=1920&screen_height=1200&browser_language=en-US&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0%20(Windows)&browser_online=true&tz_name=Asia/Qyzylorda&app_name=tiktok_web&sup_ws_ds_opt=1&version_code=270000&update_version_code=2.0.0&compress=gzip&wrss=tBx8RwpG2MJkwFrpj0zeYU6LeYcdYxfYtrOnWr1gRL0&host=https://webcast.tiktok.com&aid=1988&live_id=12&debug=false&app_language=en&room_id=7461984576673843974&identity=audience&history_comment_count=6&heartbeat_duration=0&last_rtt=397&internal_ext=fetch_time:1737383306872|start_time:0|ack_ids:,,|flag:0|seq:1|next_cursor:1737383306872_7462003955900351172_1_7462003955900350470_7461984697267002160_0|wss_info:0-1737383306872-0-0&cursor=1737383306872_7462003955900351172_1_7462003955900350470_7461984697267002160_0&history_comment_cursor=7461994261343226630&resp_content_type=protobuf&did_rule=3")
	if err != nil {
		log.Fatalf("wss url parsing err: %v", err)
	}
	vs := u.Query()
	for _, p := range wb.WsParams {
		if p.Name == "imprp" {
			continue
		}

		vs.Set(p.Name, p.Value)
	}
	vs.Set("room_id", mainURL.Query().Get("room_id"))
	u, _ = url.Parse(wb.WsUrl)
	u.RawQuery = vs.Encode()

	// Make connection to websocket
	hs := http.Header{}
	// WARN: IMPORTANT TO SEND
	hs.Add("Cookie", cookie)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), hs)
	if err != nil {
		log.Fatalf("ws dial err: %v", err)
	}

	// Ping TT server every 10 seconds
	ping := time.Tick(10 * time.Second)
	hx, err := hex.DecodeString("3A026862")
	if err != nil {
		log.Fatalf("hex decoding err: %v", err)
	}
	go func() {
		for range ping {
			err := conn.WriteMessage(websocket.BinaryMessage, hx)
			if err != nil {
				log.Fatalf("send ping err: %v", err)
			}
		}
	}()

	// Read webscoket FINALLY
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("ws read err: %v", err)
		}

		wwm := chataggr.WebcastWebsocketMessage{}
		err = proto.Unmarshal(data, &wwm)
		if err != nil {
			log.Fatalf("wwm unmarshal err: %v", err)
		}

		if wwm.Type != "msg" {
			continue
		}

		wrb := wwm.Binary

		// If has gzip header - read and unzip
		z, err := gzip.NewReader(bytes.NewReader(wrb))
		if err == nil {
			buf := &bytes.Buffer{}
			_, err = io.Copy(buf, z)
			z.Close()

			if err != nil {
				log.Fatalf("unzip copy err: %v", err)
			}

			wrb = buf.Bytes()
		}

		wr := chataggr.WebcastResponse{}
		err = proto.Unmarshal(wrb, &wr)
		if err != nil {
			log.Fatalf("wr unmarshal err: %v", err)
		}

		// Read messages and find chat
		for _, m := range wr.Messages {
			if m.Type != "WebcastChatMessage" {
				continue
			}

			wcm := chataggr.WebcastChatMessage{}
			err = proto.Unmarshal(m.Binary, &wcm)
			if err != nil {
				log.Fatalf("wr unmarshal err: %v", err)
			}

			c <- aggr.Message{
				Receiver: r,
				Text:     wcm.Comment,
				User:     wcm.User.Nickname,
			}
		}
	}
}
