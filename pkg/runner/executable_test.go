package runner_test

import (
	"context"
	"github.com/antoinetoussaint/kommence/pkg/configuration"
	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/antoinetoussaint/kommence/pkg/runner"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestSimpleExecutable(t *testing.T) {

	// Create a temporary file
	file, err := ioutil.TempFile("", "prefix")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	//var buf bytes.Buffer
	log := output.NewLogger(true)

	config := configuration.Executable{ID: "X", Cmd: "echo world", Watch: []string{file.Name()}}
	exec := runner.NewExecutable(log, &config)

	rec := make(chan output.Message, 8)
	ctx := context.Background()

	// Run in a go routine and get the messages
	go exec.Start(ctx, rec)

	// Wait a bit
	time.Sleep(500*time.Millisecond)

	// Modify the file to create a restart
	anything := []byte("anything")
	_, err = file.Write(anything)
	assert.NoError(t, err)

	// Wait a bit
	time.Sleep(500*time.Millisecond)

	// Write again
	_, err = file.Write(anything)
	assert.NoError(t, err)

	// Wait a bit
	time.Sleep(500*time.Millisecond)

	// Stop the process
	exec.Stop(ctx, rec)

	assert.Equal(t, output.Log, (<-rec).Type)
	assert.Equal(t, output.Stop, (<-rec).Type)
	assert.Equal(t, output.Restart, (<-rec).Type)
	assert.Equal(t, output.Log, (<-rec).Type)
	assert.Equal(t, output.Stop, (<-rec).Type)
	assert.Equal(t, output.Restart, (<-rec).Type)
	assert.Equal(t, output.Log, (<-rec).Type)
	assert.Equal(t, output.Stop, (<-rec).Type)

}
