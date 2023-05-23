package runner_test

import (
	"github.com/AntoineToussaint/jarvis/pkg/runner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatch(t *testing.T) {
	assert.True(t, runner.MatchPod("test", "test-9468448-j95hv"))
	assert.False(t, runner.MatchPod("test", "tester-9468448-j95hv"))
	assert.False(t, runner.MatchPod("test", "test-94684-j95hv"))
}
