package output

import "github.com/fatih/color"

const BackgroundOffset = 31

type Style = []interface{}

type Styler struct {
	current int
}

func (s *Styler) Next() Style {
	var attributes []interface{}
	if s.current < 10 {
		attributes = append(attributes, color.Attribute(BackgroundOffset+s.current))
	}
	attributes = append(attributes, color.Bold)
	s.current++
	return attributes
}
