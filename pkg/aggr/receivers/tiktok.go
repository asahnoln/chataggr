package receivers

import (
	"github.com/asahnoln/chataggr/pkg/aggr"
)

type TikTok struct{}

func NewTikTok(url string) *TikTok {
	return &TikTok{}
}

func (r *TikTok) Receive(c chan aggr.Message) {
}
