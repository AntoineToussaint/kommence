package output_test

import (
	"github.com/AntoineToussaint/jarvis/pkg/output"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseToStructured(t *testing.T) {
	l := `{"log": "info", "time": "2021-08-24T12:41:25-04:00", "msg": "test"}`
	s := output.ParseToStructured(l)
	assert.Equal(t, "INFO", s.Level)
	assert.Equal(t, "12:41:25", s.Timestamp)
	assert.Equal(t, "msg=test", s.Parsed)

	l = `{"i": "0", "time": "2021-09-08 15:24:21", "level": "info"}`
	s = output.ParseToStructured(l)
	assert.Equal(t, "INFO", s.Level)
	assert.Equal(t, "15:24:21", s.Timestamp)
	assert.Equal(t, "i=0", s.Parsed)
}
