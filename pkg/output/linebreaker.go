package output

type Message struct {
	Content string
	ID      string
}

type LineBreaker struct {
	Output      chan Message
	ID          string
	currentLine []byte
}

func NewLineBreaker(out chan Message, ID string) *LineBreaker {
	return &LineBreaker{Output: out, ID: ID}
}

func (w *LineBreaker) Write(p []byte) (int, error) {
	total := 0
	for _ , c := range p {
		total++
		if c == '\n' {
			w.Output <- Message{ID: w.ID, Content: string(w.currentLine)}
			w.currentLine = nil
			continue
		}
		w.currentLine = append(w.currentLine, c)
	}
	if w.currentLine != nil {
		w.Output <- Message{ID: w.ID, Content: string(w.currentLine)}
	}
	return total, nil
}
