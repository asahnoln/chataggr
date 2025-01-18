package receivers

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/gorilla/websocket"
)

type TikTok struct {
	url string
}

func NewTikTok(url string) *TikTok {
	return &TikTok{url}
}

func (r *TikTok) Receive(c chan aggr.Message) {
	hs := http.Header{}
	u, _ := url.Parse(r.url)
	vs := u.Query()
	// Looks constant 2nd part 6LeYcdYxfYtrOnWr1gRL0
	// 1st part takes 22 symbols:
	// zJWkjhmSlSzaZ_mWfao780
	// RWzgpDNvFwxxbxlsHtV_Ek
	// TdqOD1EBvIAJMr1jEUhpX
	vs.Set("wrss", "mvkX78gjkVSOinFk9ZTb4U6LeYcdYxfYtrOnWr1gRL0")
	u.RawQuery = vs.Encode()

	hs.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:134.0) Gecko/20100101 Firefox/134.0")
	hs.Add("Accept", "*/*")
	hs.Add("Accept-Language", "en-US,en;q=0.5")
	hs.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	// hs.Add("Sec-WebSocket-Version", "13")
	hs.Add("Origin", "https://www.tiktok.com")
	// hs.Add("Sec-WebSocket-Extensions", "permessage-deflate")
	// hs.Add("Sec-WebSocket-Key", "Da2syqtalG142CtviLX4lQ==")
	// hs.Add("Connection", "keep-alive, Upgrade")
	hs.Add("Cookie", "ttwid=1%7C5XfWiTRds5mhaylBKPImrkN_-DXgtkG-B9RRMVMYXlE%7C1735837246%7C3c8891d53acf7a5b4ad1071797338b13c30f14a4b3c6a6b20212539b912bd2a4; tt_chain_token=80E7G904W83fcHsjnaWXiA==; msToken=24lXOP4Zkd9s-0jye3F4p6e5gbS4_aNojVuErSR6y1ldjI8-Q75oe6W5ejjl2pcWrdw0WlQ33mIw9g_LwanYoY3vwnZUCZBpwi3LgyrTgbb_ImPGsyJl3BGLeDn01IzZ43AILNrbeXLykLYMjrjrCQ==; odin_tt=e926706ac66c5a2a62302f898508391000ad447a847627c7ee9262c0dc1123e16fde4a5530febd6810473e25916fd7815319b305ca339ca9b3e27c02c1fd520ffdbc05b8ca7b7a53c053bbf32b721708; passport_csrf_token=2d8e3e5d746bf7715f4a937e1da233c3; passport_csrf_token_default=2d8e3e5d746bf7715f4a937e1da233c3; multi_sids=6854322356408976389%3A39c7d85c0281664ce4fc7f1075f21568; cmpl_token=AgQQAPOsF-RO0o8hn3mF9F08_T51wNrW_4cnYNrqyQ; passport_auth_status=651f117f4ccc3d516fadff4b3d8e2715%2C; passport_auth_status_ss=651f117f4ccc3d516fadff4b3d8e2715%2C; sid_guard=39c7d85c0281664ce4fc7f1075f21568%7C1735837273%7C15552000%7CTue%2C+01-Jul-2025+17%3A01%3A13+GMT; uid_tt=eaef1ae9ffe38ee1a32b335e34a616019bdba1debe93085624e090fa9613ba63; uid_tt_ss=eaef1ae9ffe38ee1a32b335e34a616019bdba1debe93085624e090fa9613ba63; sid_tt=39c7d85c0281664ce4fc7f1075f21568; sessionid=39c7d85c0281664ce4fc7f1075f21568; sessionid_ss=39c7d85c0281664ce4fc7f1075f21568; sid_ucp_v1=1.0.0-KGEyODdhMzAxNjljNDA5Yjk0YzRjNzk2NTVkNDA4OGRiZmY1NTg3NzgKIQiFiJru4Lvcj18Q2YzbuwYYswsgDDD05P34BTgIQBJIBBADGgZtYWxpdmEiIDM5YzdkODVjMDI4MTY2NGNlNGZjN2YxMDc1ZjIxNTY4; ssid_ucp_v1=1.0.0-KGEyODdhMzAxNjljNDA5Yjk0YzRjNzk2NTVkNDA4OGRiZmY1NTg3NzgKIQiFiJru4Lvcj18Q2YzbuwYYswsgDDD05P34BTgIQBJIBBADGgZtYWxpdmEiIDM5YzdkODVjMDI4MTY2NGNlNGZjN2YxMDc1ZjIxNTY4; store-idc=maliva; store-country-code=kz; store-country-code-src=uid; tt-target-idc=useast1a; tt-target-idc-sign=hdPSGbtamqD5GvFwwuZEAoJvDzGVq5fbYOmVCRJvt69QkFjgdrIj_H5LE_h1hrs6Iz3bMyOOPseWkfEInbaSz2-gk0z6SEfvmFqdMAFYZUt2Knz-B5Bq50PnE_XeRs5pUqlll5pRKWANos5aTQTUz5EUCm4M9ZGSOi7AO3BNDtvN_h3kJhqtoy-C_OHSSMChWXELtmMpWm8Da8XHWVEV1trrayrPcF2ySN4jlkxr4m1jxUcervDFIgPUdmJ29TDU9H-aM74HRCBfR6jIgVKZ_3R7Dw9mYlhOQb7YE25A_YRZP4uMcxGWU05iYXUr8Z7hzwwUi3P9MD4f6dFnLYw_OqAHZHntAjvhic13GV6KFYu-nUWt8WufATGGAO7Y9D97TKUnEf8zKTOYTOrEitVIIUzbLup9EupF7QbGHlTaPoLuTeQyUuDdH6sGR4UodQtXm8LPtGvfhKcqovYtrGIsCPWfak9mOvS9bSK7MPDmAedDg-GWIe-9NwLeoq6QUNk7; tt_csrf_token=0vjm1tIz-p6M_LDzdB5Gmz1HEer5VfUXcKOQ; csrfToken=MKZEdlFH-3t3U6ouWGkID06FEe6_K1hQK228; s_v_web_id=verify_m5fkoyne_4Mce8XqP_d1uh_4w6l_BOTl_ucoxHn1KrPpt; ak_bmsc=D2D1A9139676B71D79BCBCF263CC7813~000000000000000000000000000000~YAAQW3kQAuR8XGqUAQAA2VsKehrUCxCHg6WM6V7OG8cEDf64uyncK64hUqZBa3D/9Og8kn3PWY3gt/I9uEqOsEja8a1jHM0HInUdbGpd8D11dt/gziSCS3fsAZet3hxCC4g/lSkIap+OjX6pak7gVXsnAlB5X/hT9i0zi8mw7/qqcDJkjJXQa4Pw7HMG8Cvn3MNYVZOHwQj6S13EpzuI1wjkwOCNltEYgqDp8Pe4G7RRXIWh/W24c9OGeANXeqPiuuL7LS/HALDiUnlyIVxfVNKYGxCTjNJL3jViWOli71X0cpqhGvnX3bgm+jYD8n+jUUfevym1dS+at3nCIH3p8cEKvAWYCy4aLdTmOZ73I6Qu6jOM1sqwZbruDM+paECyBSSGMcDgOpk=; bm_sv=B838E458B32EF4D2A889896C91E534E0~YAAQxoZlX2wRlXWUAQAAVKUbehrgLZBNgVnkNzFM/cOpkcgBIQx3q8RbdVEKqYnU1ReyVa2qpZp9ue1bTvDUWNx8hKaegiwxGVO+QNdwdb03t2C0o8atrq/Ig/YkW68hVEprRJ6I556oErASLG1X9lGJDnmH8lK4BCrAYuEjxhnbvihdzjn/oJ+fD8buUvQwnZ9L98UwDK1cRXbM5BknDwGaF058g71NjRFlzudoY3Jtv4meWUyZ0Aynn08WckH2~1")
	hs.Add("Sec-Fetch-Dest", "empty")
	hs.Add("Sec-Fetch-Mode", "websocket")
	hs.Add("Sec-Fetch-Site", "same-site")
	hs.Add("Pragma", "no-cache")
	hs.Add("Cache-Control", "no-cache")
	// hs.Add("Upgrade", "websocket")

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), hs)
	if err != nil {
		log.Printf("ws resp: %+v", resp)
		log.Fatalf("ws dial err: %v", err)
	}
	log.Print("received conn")
	b, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	log.Printf("ws resp: %s", b)

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ws read err: %v", err)
			return
		}

		log.Printf("tiktok msg: %s", data)

		// c <- aggr.Message{
		// 	Text: text,
		// 	User: findSubstrBetween(msg, "display-name=", ";"),
		// }
	}

	c <- aggr.Message{User: "someone", Text: "Hi Tik"}
	c <- aggr.Message{}
}
