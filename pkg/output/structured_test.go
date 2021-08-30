package output_test

import (
	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseToStructured(t *testing.T) {
	l := `{"log": "info", "time": "2021-08-24T12:41:25-04:00", "msg": "test"}`
	s := output.ParseToStructured(l)
	assert.Equal(t, "info", s.Level)
	assert.False(t, s.Timestamp.IsZero())
	assert.Equal(t, "msg=test", s.Parsed)
}
