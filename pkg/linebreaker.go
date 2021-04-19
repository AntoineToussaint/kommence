package pkg

type LineBreaker struct {
	Output      chan Message
	currentLine []byte
}

func NewLineBreaker(out chan Message) *LineBreaker {
	return &LineBreaker{Output: out}
}

func (w *LineBreaker) Write(p []byte) (int, error) {
	total := 0
	for _ , c := range p {
		total++
		if c == '\n' {
			w.Output <- Message{Content: string(w.currentLine)}
			w.currentLine = nil
			continue
		}
		w.currentLine = append(w.currentLine, c)
	}
	if w.currentLine != nil {
		w.Output <- Message{Content: string(w.currentLine)}
	}
	return total, nil
}
