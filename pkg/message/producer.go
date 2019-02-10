package message

// Producer produces messages to messaging system.
type Producer interface {
	Open()
	Produce(msg []byte) error
	Close()
}
