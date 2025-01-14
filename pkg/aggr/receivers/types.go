package receivers

type Connector interface {
	Connect(chan []byte)
}
