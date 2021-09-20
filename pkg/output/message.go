package output

type MessageType int

const (
	Log MessageType = iota
	Error
	Stop
	Restart
	PodConnection
	Memory
	CPU
)

// Message are how processes communicate
type Message struct {
	// ID is the message sender
	ID string
	// Type of message
	Type MessageType
	// Content of the message
	Content string
}
