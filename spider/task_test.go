package spider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	err := addTaskCtrl(1, cancel)
	assert.Nil(t, err)
	err = addTaskCtrl(1, cancel)
	assert.NotNil(t, err)

	ok := CancelTask(1)
	assert.True(t, ok)
	ok = CancelTask(1)
	assert.False(t, ok)
}
