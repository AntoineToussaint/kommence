package output

type Message struct {
	Content string
	ID      string
	IsError bool
}

type LineBreaker struct {
	Output      chan Message
	ID          string
	IsError bool
	currentLine []byte
}

func NewLineBreaker(out chan Message, ID string, isError bool) *LineBreaker {
	return &LineBreaker{Output: out, ID: ID, IsError: isError}
}

func (w *LineBreaker) Write(p []byte) (int, error) {
	total := 0
	for _ , c := range p {
		total++
		if c == '\n' {
			w.Output <- Message{ID: w.ID, Content: string(w.currentLine), IsError: w.IsError}
			w.currentLine = nil
			continue
		}
		w.currentLine = append(w.currentLine, c)
	}
	if w.currentLine != nil {
		w.Output <- Message{ID: w.ID, Content: string(w.currentLine), IsError: w.IsError}
	}
	return total, nil
}
