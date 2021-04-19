package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"strings"
)

type FormatterConfiguration struct {
	Color string
	Json  string
}

type Formatter struct {
	FormatterConfiguration
	source Source
	padding string
}

func (d Formatter) ID() string {
	return d.source.ID()
}

func (d Formatter) Start(ctx context.Context) {
	d.source.Start(ctx)
}

func NewFormatter(c FormatterConfiguration, source Source, maxLength int) (*Formatter, error) {
	pad := maxLength - len(source.ID())
	padding := ""
	for i := 0 ; i<pad ;i++ {
		padding += " "
	}
	return &Formatter{FormatterConfiguration: c, source: source, padding: padding}, nil
}

type flat = map[string]interface{}

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
	out.Content = fmt.Sprintf("%v(%v) %v", d.padding, d.source.ID(), out.Content)
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
