package tests

import (
	"github.com/sudachen/go-numerai/numerai"
	"gotest.tools/assert"
	"testing"
)

func Test_Rounds(t *testing.T) {
	curr := numerai.CurrentRound()
	assert.Assert(t, curr.Status == numerai.OPEN)
}
