package output_test

import (
	"bytes"
	"github.com/AntoineToussaint/jarvis/pkg/output"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	log := output.NewLogger(false, output.WithOut(&buf))
	log.Printf("hello: %v", "world")
	assert.Equal(t, "hello: world", buf.String())
}
