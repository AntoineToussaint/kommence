package pkg

type LineBreaker struct {
	Output      chan string
	currentLine []byte
}

func NewLineBreaker(out chan string) *LineBreaker {
	return &LineBreaker{Output: out}
}

func (w *LineBreaker) Write(p []byte) (int, error) {
	total := 0
	for _ , c := range p {
		total++
		if c == '\n' {
			w.Output <- string(w.currentLine)
			w.currentLine = []byte{}
			continue
		}
		w.currentLine = append(w.currentLine, c)
	}
	return total, nil
}


// Close flushes the last of the output into the underlying writer.
func (w *LineBreaker) Close() error {
	w.Output <- string(w.currentLine)
	return nil
}
