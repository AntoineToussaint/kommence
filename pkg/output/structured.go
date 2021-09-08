package output

import (
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"strings"
	"time"
)

const TimeFormat = "15:04:05"

type StructuredLog struct {
	Level string
	Timestamp string
	Parsed string
}

type Data map[string]string

func matchLevel(s string) (string, bool) {
	switch strings.ToLower(s) {
	case "info":
		return "INFO", true
	case "debug":
		return "DEBUG", true
	case "error":
		return "ERROR", true
	case "warn":
		return "WARN", true
	}
	return "", false
}

func matchTimestamp(s string) (time.Time, bool) {
	t, err := dateparse.ParseLocal(s)
	return t, err == nil
}

func ParseToStructured(l string) StructuredLog {
	s := StructuredLog{
		Parsed: l,
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
		if lvl, ok := matchLevel(v); ok {
			strip = append(strip, k)
			s.Level = lvl
		}
		if t, ok := matchTimestamp(v); ok {
			strip = append(strip, k)
			s.Timestamp = t.Format(TimeFormat)
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
