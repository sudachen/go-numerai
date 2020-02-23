package tests

import (
	"github.com/sudachen/go-numerai/nuai"
	"gotest.tools/assert"
	"testing"
)

func Test_Rounds(t *testing.T) {
	n := &nuai.Nuai{}
	curr := n.CurrentRound()
	assert.Assert(t, curr.Status == nuai.OPEN)
}
