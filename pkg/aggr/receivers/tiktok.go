package receivers

import (
	"github.com/Davincible/gotiktoklive"
	"github.com/asahnoln/chataggr/pkg/aggr"
)

type TikTok struct{ l *gotiktoklive.Live }

func NewTikTok(l *gotiktoklive.Live) *TikTok {
	return &TikTok{l}
}

func (r *TikTok) Receive(c chan aggr.Message) {
	for event := range r.l.Events {
		switch event.(type) {
		case gotiktoklive.ChatEvent:
			c <- aggr.Message{}
		}
	}
}
