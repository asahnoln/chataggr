package streamer

// Message contains text from a chat
type Message struct{}

type Messager interface {
	Receive(c chan Message)
}

func Run(mrs []Messager, c chan Message) {
	for _, m := range mrs {
		go m.Receive(c)
	}
}
