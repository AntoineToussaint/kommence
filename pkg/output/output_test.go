package output_test

import (
	"bytes"
	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	log := output.NewLogger(false, output.WithOut(&buf))
	log.Printf("hello: %v", "world")
	assert.Equal(t, buf.String(), "hello: world")
}
