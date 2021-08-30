package output

import (
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"strings"
	"time"
)

type StructuredLog struct {
	Raw string
	Level string
	Timestamp time.Time
	Parsed string
}

type Data map[string]string

func matchLevel(s string) bool {
	switch s {
	case "info":
		return true
	case "INFO":
		return true
	case "debug":
		return true
	case "DEBUG":
		return true
	case "error":
		return true
	case "ERROR":
		return true
	case "warn":
		return true
	case "WARN":
		return true
	}
	return false
}

func matchTimestamp(s string) (time.Time, bool) {
	t, err := dateparse.ParseLocal(s)
	return t, err == nil
}

func ParseToStructured(l string) StructuredLog {
	s := StructuredLog{
		Raw: l,
	}
	// Try to parse JSON parse
	var data Data
	err := json.Unmarshal([]byte(l), &data)
	if err != nil {
		return s
	}
	var strip []string
	// Find a Level
	for k, v := range data {
		if matchLevel(v) {
			strip = append(strip, k)
			s.Level = v
		}
		if t, ok := matchTimestamp(v); ok {
			strip = append(strip, k)
			s.Timestamp = t
		}
	}
	for _, c := range strip {
		delete(data, c)
	}
	var msgs []string
	for k, v := range data {
		msgs = append(msgs, fmt.Sprintf("%v=%v", k, v))
	}
	s.Parsed = strings.Join(msgs, " ")
	return s
}
