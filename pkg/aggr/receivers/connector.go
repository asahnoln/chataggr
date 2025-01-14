package receivers

type Connector interface {
	Connect(chan []byte)
}

type TwitchConnector struct{}

func NewTwitchConnector() *TwitchConnector { return nil }

func (conn *TwitchConnector) Connect(c chan []byte) {
}
