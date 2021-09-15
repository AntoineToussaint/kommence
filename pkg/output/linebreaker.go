package output

type LineBreaker struct {
	Output      chan Message
	ID          string
	messageType MessageType
	currentLine []byte
}

func NewLineBreaker(out chan Message, ID string, t MessageType) *LineBreaker {
	return &LineBreaker{Output: out, ID: ID, messageType: t}
}

func (w *LineBreaker) Write(p []byte) (int, error) {
	total := 0
	for _, c := range p {
		total++
		if c == '\n' {
			w.Output <- Message{ID: w.ID, Type: w.messageType, Content: string(w.currentLine)}
			w.currentLine = nil
			continue
		}
		w.currentLine = append(w.currentLine, c)
	}
	if w.currentLine != nil {
		w.Output <- Message{ID: w.ID, Type: w.messageType, Content: string(w.currentLine)}
	}
	return total, nil
}
