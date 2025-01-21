// package aggr is for launching all message receivers
// and sending data to a channel
package aggr

// Message contains text from a chat
type Message struct {
	User, Text string
	Receiver   Receiver
}

// Receiver is a chat service.
// Receiver should parse messages from a chat
// and send them to the given channel
// when they are ready.
type Receiver interface {
	Receive(c chan Message)
}

// Run launches all receivers in goroutines
// so that they connect to their chats
// and send parsed messages to the given channel
func Run(rs []Receiver, c chan Message) {
	for _, r := range rs {
		go r.Receive(c)
	}
}
