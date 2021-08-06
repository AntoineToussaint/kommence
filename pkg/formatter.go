package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"strings"
)

type Type int

const (
	RunType Type = iota
	KubeType
)

type FormatterConfiguration struct {
	Color string
	Json  string
}

type Formatter struct {
	FormatterConfiguration
	source  Source
	padding string
	Type    Type
}

func (d Formatter) ID() string {
	return d.source.ID()
}

func (d Formatter) Start(ctx context.Context) {
	d.source.Start(ctx)
}

var col int

func pickColor() string {
	colors := []string{"red", "blue", "yellow", "green"}
	col++
	return colors[col % 4]
}

func NewFormatter(t Type, c FormatterConfiguration, source Source, maxLength int) (*Formatter, error) {
	pad := maxLength - len(source.ID())
	padding := ""
	for i := 0 ; i<pad ;i++ {
		padding += " "
	}
	if c.Color == "" {
		c.Color = pickColor()
	}
	return &Formatter{Type: t, FormatterConfiguration: c, source: source, padding: padding}, nil
}

type flat = map[string]interface{}

var prefixes map[Type]string

func init() {
	prefixes = make(map[Type]string)
	prefixes[RunType] = "Ⓡ"
	prefixes[KubeType] = "⎈"
}

func (d Formatter) Format(msg Message) Message {
	out := Message{Content: msg.Content}
	if d.Json == "flatten" {
		var f flat
		err := json.Unmarshal([]byte(msg.Content), &f)
		if err == nil {
			var flat []string
			for k, v := range f {
				flat = append(flat, fmt.Sprintf("%v=%v", k, v))
			}
			out.Content = strings.Join(flat, " ")
		}
	}
	out.Content = fmt.Sprintf("(%v %v)%v %v", prefixes[d.Type], d.source.ID(), d.padding, out.Content)
	switch d.Color {
	case "red":
		out.Content = color.FgRed.Render(out.Content)
	case "blue":
		out.Content = color.FgBlue.Render(out.Content)

	}
	return out
}

func (d Formatter) Produce(ctx context.Context) <-chan Message {
	out := make(chan Message)
	go func() {
		for msg := range d.source.Produce(ctx) {
			out <- d.Format(msg)
		}
	}()
	return out
}
