package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/asahnoln/chataggr/pkg/aggr"
	"github.com/asahnoln/chataggr/pkg/aggr/receivers"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// stubInMemoryReceiver is used for debugging purposes
type stubInMemoryReceiver struct{}

func (r *stubInMemoryReceiver) Receive(c chan aggr.Message) {
	t := time.Tick(3 * time.Second)
	for range t {
		c <- aggr.Message{User: "asahnoln", Text: "Hi"}
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c := make(chan aggr.Message)
	aggr.Run([]aggr.Receiver{
		createTwitchReceiver(),
		createTikTokReceiver(),
	}, c)

	// TODO: Suppress info from rl
	s, closeAudioDevice, unloadSound := initSound()
	defer closeAudioDevice()
	defer unloadSound(s)

loop:
	for {
		select {
		case m := <-c:
			printMessage(m)
			playSound(s)
		case <-interrupt:
			break loop
		}
	}
}

func initSound() (rl.Sound, func(), func(rl.Sound)) {
	rl.InitAudioDevice()

	s := rl.LoadSound(os.Getenv("ALERT_SND_PATH"))
	if !rl.IsSoundValid(s) {
		log.Fatalf("Failed to load sound: %v\n", s)
	}

	return s, rl.CloseAudioDevice, rl.UnloadSound
}

func playSound(s rl.Sound) {
	rl.PlaySound(s)
}

func printMessage(m aggr.Message) {
	r := getReceiverName(m)
	fmt.Printf("%s:\t%s\t[%s]: %s\n", r, m.User, time.Now().Format(time.TimeOnly), m.Text)
}

func getReceiverName(m aggr.Message) string {
	r := "?"
	switch m.Receiver.(type) {
	case *receivers.Twitch:
		r = "TW"
	case *receivers.TikTok:
		r = "TT"
	}
	return r
}

func createTikTokReceiver() *receivers.TikTok {
	return receivers.NewTikTok("https://webcast.tiktok.com/webcast/im/fetch/?version_code=180800&device_platform=web&cookie_enabled=true&screen_width=1920&screen_height=1200&browser_language=en-US&browser_platform=Win32&browser_name=Mozilla&browser_version=5.0%20(Windows)&browser_online=true&tz_name=Asia/Qyzylorda&aid=1988&app_name=tiktok_web&live_id=12&version_code=270000&debug=false&app_language=en&room_id=7462709674645605125&identity=audience&history_comment_count=6&fetch_rule=1&last_rtt=-1&cursor=0&internal_ext=0&sup_ws_ds_opt=1&resp_content_type=protobuf&did_rule=3&msToken=FMOmN2OX9D6lTwDva88JEX1cfiurz-tYT-PDm9tucTtkzaPa37Igo5TT_bhEM4SyAC_1KrE9-OPCi53Evos1V23Nt743uG0i9NBx7E-sxG_ytGNAolT4DXXT9R78EdrZVqFajpe4etcuY4dd-bSbuw==&X-Bogus=DFSzsIVOtUxANa7Utdf-wn2J46B8&_signature=_02B4Z6wo00001uSstHQAAIDBCxtzj-ndYr7krrDAAN7J33")
}

func createTwitchReceiver() *receivers.Twitch {
	return receivers.NewTwitch("wss://irc-ws.chat.twitch.tv/")
}

func createStubInMemoryReceiver() *stubInMemoryReceiver {
	return &stubInMemoryReceiver{}
}
