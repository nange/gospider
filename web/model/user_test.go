package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenUserHashPassword(t *testing.T) {
	hash, err := GenUserHashPassword("admin")
	assert.NoErrorf(t, err, "shoudl gen hash password success")
	t.Logf("hash:%s\n", hash)
}
