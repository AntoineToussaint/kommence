package output

import (
	"github.com/fatih/color"
	"io"
	"os"
)

type Logger struct {
	out   io.Writer
	err   io.Writer
	debug bool
}

type LoggerOption = func(*Logger)

func WithOut(out io.Writer) LoggerOption {
	return func(logger *Logger) {
		logger.out = out
	}
}

func NewLogger(debug bool, opts ...LoggerOption) *Logger {
	logger := &Logger{
		out:   os.Stdout,
		err:   os.Stderr,
		debug: debug,
	}
	for _, opt := range opts {
		opt(logger)
	}
	return logger
}

func (l *Logger) Printf(s string, args ...interface{}) {
	var attributes []color.Attribute
	var fmtArgs []interface{}
	for _, arg := range args {
		if attr, ok := arg.(color.Attribute); ok {
			attributes = append(attributes, attr)
			continue
		}
		fmtArgs = append(fmtArgs, arg)
	}
	c := color.New(attributes...)
	msg := c.Sprintf(s, fmtArgs...)
	_, _ = l.out.Write([]byte(msg))
}

func (l *Logger) Debugf(s string, args ...interface{}) {
	if l.debug {
		var attributes []color.Attribute
		var fmtArgs []interface{}
		for _, arg := range args {
			if attr, ok := arg.(color.Attribute); ok {
				attributes = append(attributes, attr)
				continue
			}
			fmtArgs = append(fmtArgs, arg)
		}
		c := color.New(attributes...)
		msg := c.Sprintf("[DEBUG] "+s, fmtArgs...)
		_, _ = l.out.Write([]byte(msg))
	}
}
func (l *Logger) Errorf(s string, args ...interface{}) {
	var attributes []color.Attribute
	var fmtArgs []interface{}
	for _, arg := range args {
		if attr, ok := arg.(color.Attribute); ok {
			attributes = append(attributes, attr)
			continue
		}
		fmtArgs = append(fmtArgs, arg)
	}
	c := color.New(attributes...)
	msg := c.Sprintf("[ERROR] "+s, fmtArgs...)
	_, _ = l.err.Write([]byte(msg))
}
