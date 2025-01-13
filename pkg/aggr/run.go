package aggr

// Message contains text from a chat
type Message struct{}

type Receiver interface {
	Receive(c chan Message)
}

func Run(mrs []Receiver, c chan Message) {
	for _, m := range mrs {
		go m.Receive(c)
	}
}
