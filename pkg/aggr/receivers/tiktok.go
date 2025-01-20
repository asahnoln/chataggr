package receivers

import (
	"encoding/hex"
	"fmt"
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

func (r *TikTok) Receive(c chan aggr.Message) {
	// TODO: fetch room id
	// TODO: fetch msToken
	// TODO: fetch X-Bogus
	// TODO: make signature
	mainUrl, err := url.Parse(r.url)
	if err != nil {
		log.Fatalf("parsing webcast url err: %v", err)
	}
	req, err := http.NewRequest(http.MethodGet, mainUrl.String(), nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:134.0) Gecko/20100101 Firefox/134.0")
	// req.Header.Add("Accept", "*/*")
	// req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	// // req.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	// req.Header.Add("Referer", "https://www.tiktok.com/")
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	// req.Header.Add("Origin", "https://www.tiktok.com")
	// req.Header.Add("Sec-Fetch-Dest", "empty")
	// req.Header.Add("Sec-Fetch-Mode", "cors")
	// req.Header.Add("Sec-Fetch-Site", "same-site")
	// req.Header.Add("Connection", "keep-alive")
	// TODO: fetch cookie?
	// WARN: IMPORTANT TO SEND
	req.Header.Add("Cookie", "ttwid=1%7C5XfWiTRds5mhaylBKPImrkN_-DXgtkG-B9RRMVMYXlE%7C1735837246%7C3c8891d53acf7a5b4ad1071797338b13c30f14a4b3c6a6b20212539b912bd2a4; tt_chain_token=80E7G904W83fcHsjnaWXiA==; msToken=uucqJn5zUgf2J8wHZkGkRxGAcksKLJxnQWFyQwV6s6hOx6wPPcwDkVfC9InnSGY5tDzJmQCzUtvrP3IiClXCZGQz_3W4KLvEEbKdqREHgwtQYPVy2U_Ss48Ty6VaU55rS2QnYouUGNYPVCftqz-oQA==; odin_tt=e926706ac66c5a2a62302f898508391000ad447a847627c7ee9262c0dc1123e16fde4a5530febd6810473e25916fd7815319b305ca339ca9b3e27c02c1fd520ffdbc05b8ca7b7a53c053bbf32b721708; passport_csrf_token=2d8e3e5d746bf7715f4a937e1da233c3; passport_csrf_token_default=2d8e3e5d746bf7715f4a937e1da233c3; multi_sids=6854322356408976389%3A39c7d85c0281664ce4fc7f1075f21568; cmpl_token=AgQQAPOsF-RO0o8hn3mF9F08_T51wNrW_4cnYNrqyQ; passport_auth_status=651f117f4ccc3d516fadff4b3d8e2715%2C; passport_auth_status_ss=651f117f4ccc3d516fadff4b3d8e2715%2C; sid_guard=39c7d85c0281664ce4fc7f1075f21568%7C1735837273%7C15552000%7CTue%2C+01-Jul-2025+17%3A01%3A13+GMT; uid_tt=eaef1ae9ffe38ee1a32b335e34a616019bdba1debe93085624e090fa9613ba63; uid_tt_ss=eaef1ae9ffe38ee1a32b335e34a616019bdba1debe93085624e090fa9613ba63; sid_tt=39c7d85c0281664ce4fc7f1075f21568; sessionid=39c7d85c0281664ce4fc7f1075f21568; sessionid_ss=39c7d85c0281664ce4fc7f1075f21568; sid_ucp_v1=1.0.0-KGEyODdhMzAxNjljNDA5Yjk0YzRjNzk2NTVkNDA4OGRiZmY1NTg3NzgKIQiFiJru4Lvcj18Q2YzbuwYYswsgDDD05P34BTgIQBJIBBADGgZtYWxpdmEiIDM5YzdkODVjMDI4MTY2NGNlNGZjN2YxMDc1ZjIxNTY4; ssid_ucp_v1=1.0.0-KGEyODdhMzAxNjljNDA5Yjk0YzRjNzk2NTVkNDA4OGRiZmY1NTg3NzgKIQiFiJru4Lvcj18Q2YzbuwYYswsgDDD05P34BTgIQBJIBBADGgZtYWxpdmEiIDM5YzdkODVjMDI4MTY2NGNlNGZjN2YxMDc1ZjIxNTY4; store-idc=maliva; store-country-code=kz; store-country-code-src=uid; tt-target-idc=useast1a; tt-target-idc-sign=hdPSGbtamqD5GvFwwuZEAoJvDzGVq5fbYOmVCRJvt69QkFjgdrIj_H5LE_h1hrs6Iz3bMyOOPseWkfEInbaSz2-gk0z6SEfvmFqdMAFYZUt2Knz-B5Bq50PnE_XeRs5pUqlll5pRKWANos5aTQTUz5EUCm4M9ZGSOi7AO3BNDtvN_h3kJhqtoy-C_OHSSMChWXELtmMpWm8Da8XHWVEV1trrayrPcF2ySN4jlkxr4m1jxUcervDFIgPUdmJ29TDU9H-aM74HRCBfR6jIgVKZ_3R7Dw9mYlhOQb7YE25A_YRZP4uMcxGWU05iYXUr8Z7hzwwUi3P9MD4f6dFnLYw_OqAHZHntAjvhic13GV6KFYu-nUWt8WufATGGAO7Y9D97TKUnEf8zKTOYTOrEitVIIUzbLup9EupF7QbGHlTaPoLuTeQyUuDdH6sGR4UodQtXm8LPtGvfhKcqovYtrGIsCPWfak9mOvS9bSK7MPDmAedDg-GWIe-9NwLeoq6QUNk7; tt_csrf_token=0vjm1tIz-p6M_LDzdB5Gmz1HEer5VfUXcKOQ; csrfToken=MKZEdlFH-3t3U6ouWGkID06FEe6_K1hQK228; s_v_web_id=verify_m5fkoyne_4Mce8XqP_d1uh_4w6l_BOTl_ucoxHn1KrPpt; csrf_session_id=319e32e5d82cbdfecd6c1e62ac4f67e3; ak_bmsc=1943FF6296979185979153E7E6477AC0~000000000000000000000000000000~YAAQiNhUuGk+KGGUAQAAmgsbhBoy9fXWrl87a5/LuseKH2gA72I3YVJH2wpIW4/v9mm6bXMNcnG8rw+khE/8BzbsYoj28NaFbHRhyiRXreDrTvyhT8jwJ3fOoOnsPUBwMTHb4sw677DvnNQlTdnw2VX7sREsBWfsuIEAMEQPyhpOn5knBjp/YglkBsN6FjCgcMckNgbf/qS3ivIKYqLKQVXbmE5lReCmQKA/gk12khOrJy851UOOAW+wUD0BoO14imjxUZz6GwLrDFXA+KinlDjnrwRKUxX4nFB/gEKovulA0tQUf9qFhQlvK9kRcsTUyKEEwXGvRgjydBGidjh89R95zwlMmuNiyzE7N2xgMJ0pBcFFU904JPR/BpVvEGA47AReBvOV+cls; bm_sv=E85210612B9364AE5D0E360A370C84A1~YAAQiNhUuJ5GKGGUAQAATlgdhBoc9yx57Q3wNI+nxUgQeHuTpZ3GQyx8fjm+yS0h/YjMRWUKFAMacrabe8B1yYIKhAExX0+8hAV3UWQwSxiYJfUHCt6HOCR3Py3Aex94R7L45J0NZpVa2xA8rEkfr4tlRr/QwdAu4vq7EygWJuikCVSiFNVEMgJ6NDHHnxxxMVkd6/j3iCCL/mXcdWxa72RFEWcP6JfmxdYL90SjGIgs0o2QeuQanMe3tCwM2EWf~1")
	// req.Header.Add("TE", "trailer")

	if err != nil {
		log.Fatalf("webcast req err: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("webcast do err: %v", err)
	}
	// TODO: Parse protobuf for wrss param
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("webcast body read err: %v", err)
	}

	var wb chataggr.WebcastResponse
	if err := proto.Unmarshal(b, &wb); err != nil {
		log.Fatalf("unmarshal proto err: %v", err)
	}
	log.Printf("webcast ws url and wrss: %v, %+v", wb.WsUrl, wb.WsParams)

	hs := http.Header{}
	// WARN: For now taking URL straight from browser
	u, _ := url.Parse("wss://webcast16-ws-useast1a.tiktok.com/webcast/im/ws_proxy/ws_reuse_supplement/?version_code=180800&device_platform=web&cookie_enabled=true&screen_width=1920&screen_height=1200&browser_language=en-US&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0%20(Windows)&browser_online=true&tz_name=Asia/Qyzylorda&app_name=tiktok_web&sup_ws_ds_opt=1&version_code=270000&update_version_code=2.0.0&compress=gzip&wrss=tBx8RwpG2MJkwFrpj0zeYU6LeYcdYxfYtrOnWr1gRL0&host=https://webcast.tiktok.com&aid=1988&live_id=12&debug=false&app_language=en&room_id=7461984576673843974&identity=audience&history_comment_count=6&heartbeat_duration=0&last_rtt=397&internal_ext=fetch_time:1737383306872|start_time:0|ack_ids:,,|flag:0|seq:1|next_cursor:1737383306872_7462003955900351172_1_7462003955900350470_7461984697267002160_0|wss_info:0-1737383306872-0-0&cursor=1737383306872_7462003955900351172_1_7462003955900350470_7461984697267002160_0&history_comment_cursor=7461994261343226630&resp_content_type=protobuf&did_rule=3")
	vs := u.Query()
	for _, p := range wb.WsParams {
		if p.Name == "imprp" {
			continue
		}

		vs.Set(p.Name, p.Value)
	}
	u, _ = url.Parse(wb.WsUrl)
	// Looks constant 2nd part 6LeYcdYxfYtrOnWr1gRL0, seems to be dependent on User-Agent
	// 1st part takes 22 symbols:
	// zJWkjhmSlSzaZ_mWfao780
	// RWzgpDNvFwxxbxlsHtV_Ek
	// vs.Set("wrss", "mvkX78gjkVSOinFk9ZTb4U6LeYcdYxfYtrOnWr1gRL0")
	u.RawQuery = vs.Encode()

	// hs.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:134.0) Gecko/20100101 Firefox/134.0")
	// hs.Add("Accept", "*/*")
	// hs.Add("Accept-Language", "en-US,en;q=0.5")
	// hs.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	// hs.Add("Origin", "https://www.tiktok.com")

	// WARN: IMPORTANT TO SEND
	hs.Add("Cookie", "ttwid=1%7C5XfWiTRds5mhaylBKPImrkN_-DXgtkG-B9RRMVMYXlE%7C1735837246%7C3c8891d53acf7a5b4ad1071797338b13c30f14a4b3c6a6b20212539b912bd2a4; tt_chain_token=80E7G904W83fcHsjnaWXiA==; msToken=nvpse9Ks1b7pbcK5LS7k3hSoPaduiTYDJwT4y8z6zogy4oU2IfjJ04pMCNK5_OaN5_fgUzxu4IY5f31vh9E4ugF-iISgjVWrQMAxZDPQp9wlyqs6ErgC2xbl9aZM_j_IqLMpoSg4x4K8kMcSceRNFg==; odin_tt=e926706ac66c5a2a62302f898508391000ad447a847627c7ee9262c0dc1123e16fde4a5530febd6810473e25916fd7815319b305ca339ca9b3e27c02c1fd520ffdbc05b8ca7b7a53c053bbf32b721708; passport_csrf_token=2d8e3e5d746bf7715f4a937e1da233c3; passport_csrf_token_default=2d8e3e5d746bf7715f4a937e1da233c3; multi_sids=6854322356408976389%3A39c7d85c0281664ce4fc7f1075f21568; cmpl_token=AgQQAPOsF-RO0o8hn3mF9F08_T51wNrW_4cnYNrqyQ; passport_auth_status=651f117f4ccc3d516fadff4b3d8e2715%2C; passport_auth_status_ss=651f117f4ccc3d516fadff4b3d8e2715%2C; sid_guard=39c7d85c0281664ce4fc7f1075f21568%7C1735837273%7C15552000%7CTue%2C+01-Jul-2025+17%3A01%3A13+GMT; uid_tt=eaef1ae9ffe38ee1a32b335e34a616019bdba1debe93085624e090fa9613ba63; uid_tt_ss=eaef1ae9ffe38ee1a32b335e34a616019bdba1debe93085624e090fa9613ba63; sid_tt=39c7d85c0281664ce4fc7f1075f21568; sessionid=39c7d85c0281664ce4fc7f1075f21568; sessionid_ss=39c7d85c0281664ce4fc7f1075f21568; sid_ucp_v1=1.0.0-KGEyODdhMzAxNjljNDA5Yjk0YzRjNzk2NTVkNDA4OGRiZmY1NTg3NzgKIQiFiJru4Lvcj18Q2YzbuwYYswsgDDD05P34BTgIQBJIBBADGgZtYWxpdmEiIDM5YzdkODVjMDI4MTY2NGNlNGZjN2YxMDc1ZjIxNTY4; ssid_ucp_v1=1.0.0-KGEyODdhMzAxNjljNDA5Yjk0YzRjNzk2NTVkNDA4OGRiZmY1NTg3NzgKIQiFiJru4Lvcj18Q2YzbuwYYswsgDDD05P34BTgIQBJIBBADGgZtYWxpdmEiIDM5YzdkODVjMDI4MTY2NGNlNGZjN2YxMDc1ZjIxNTY4; store-idc=maliva; store-country-code=kz; store-country-code-src=uid; tt-target-idc=useast1a; tt-target-idc-sign=hdPSGbtamqD5GvFwwuZEAoJvDzGVq5fbYOmVCRJvt69QkFjgdrIj_H5LE_h1hrs6Iz3bMyOOPseWkfEInbaSz2-gk0z6SEfvmFqdMAFYZUt2Knz-B5Bq50PnE_XeRs5pUqlll5pRKWANos5aTQTUz5EUCm4M9ZGSOi7AO3BNDtvN_h3kJhqtoy-C_OHSSMChWXELtmMpWm8Da8XHWVEV1trrayrPcF2ySN4jlkxr4m1jxUcervDFIgPUdmJ29TDU9H-aM74HRCBfR6jIgVKZ_3R7Dw9mYlhOQb7YE25A_YRZP4uMcxGWU05iYXUr8Z7hzwwUi3P9MD4f6dFnLYw_OqAHZHntAjvhic13GV6KFYu-nUWt8WufATGGAO7Y9D97TKUnEf8zKTOYTOrEitVIIUzbLup9EupF7QbGHlTaPoLuTeQyUuDdH6sGR4UodQtXm8LPtGvfhKcqovYtrGIsCPWfak9mOvS9bSK7MPDmAedDg-GWIe-9NwLeoq6QUNk7; tt_csrf_token=0vjm1tIz-p6M_LDzdB5Gmz1HEer5VfUXcKOQ; csrfToken=MKZEdlFH-3t3U6ouWGkID06FEe6_K1hQK228; s_v_web_id=verify_m5fkoyne_4Mce8XqP_d1uh_4w6l_BOTl_ucoxHn1KrPpt; ak_bmsc=1943FF6296979185979153E7E6477AC0~000000000000000000000000000000~YAAQiNhUuGk+KGGUAQAAmgsbhBoy9fXWrl87a5/LuseKH2gA72I3YVJH2wpIW4/v9mm6bXMNcnG8rw+khE/8BzbsYoj28NaFbHRhyiRXreDrTvyhT8jwJ3fOoOnsPUBwMTHb4sw677DvnNQlTdnw2VX7sREsBWfsuIEAMEQPyhpOn5knBjp/YglkBsN6FjCgcMckNgbf/qS3ivIKYqLKQVXbmE5lReCmQKA/gk12khOrJy851UOOAW+wUD0BoO14imjxUZz6GwLrDFXA+KinlDjnrwRKUxX4nFB/gEKovulA0tQUf9qFhQlvK9kRcsTUyKEEwXGvRgjydBGidjh89R95zwlMmuNiyzE7N2xgMJ0pBcFFU904JPR/BpVvEGA47AReBvOV+cls; bm_sv=E85210612B9364AE5D0E360A370C84A1~YAAQiNhUuJ5GKGGUAQAATlgdhBoc9yx57Q3wNI+nxUgQeHuTpZ3GQyx8fjm+yS0h/YjMRWUKFAMacrabe8B1yYIKhAExX0+8hAV3UWQwSxiYJfUHCt6HOCR3Py3Aex94R7L45J0NZpVa2xA8rEkfr4tlRr/QwdAu4vq7EygWJuikCVSiFNVEMgJ6NDHHnxxxMVkd6/j3iCCL/mXcdWxa72RFEWcP6JfmxdYL90SjGIgs0o2QeuQanMe3tCwM2EWf~1")

	// hs.Add("Sec-Fetch-Dest", "empty")
	// hs.Add("Sec-Fetch-Mode", "websocket")
	// hs.Add("Sec-Fetch-Site", "same-site")
	// hs.Add("Pragma", "no-cache")
	// hs.Add("Cache-Control", "no-cache")
	//
	// hs.Add("Sec-WebSocket-Version", "13")
	// hs.Add("Sec-WebSocket-Extensions", "permessage-deflate")
	// hs.Add("Sec-WebSocket-Key", "Da2syqtalG142CtviLX4lQ==")
	// hs.Add("Connection", "keep-alive, Upgrade")
	// hs.Add("Upgrade", "websocket")

	log.Printf("ws url: %v", u.String())
	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), hs)
	if err != nil {
		log.Printf("ws resp: %+v", resp)
		log.Fatalf("ws dial err: %v", err)
	}
	b, _ = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	log.Printf("ws resp: %s", b)

	ping := time.Tick(10 * time.Second)
	hx, err := hex.DecodeString("3A026862")
	if err != nil {
		log.Fatalf("hex decoding err: %v", err)
	}
	go func() {
		for range ping {
			conn.WriteMessage(websocket.BinaryMessage, hx)
		}
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ws read err: %v", err)
			return
		}

		fmt.Printf("tiktok msg: %s\n", data)
		log.Printf("tiktok msg: %s", data)

		c <- aggr.Message{
			Text: string(data),
		}
	}

	// c <- aggr.Message{User: "someone", Text: "Hi Tik"}
	// c <- aggr.Message{}
}
