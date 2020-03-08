package tests

import (
	"github.com/sudachen/go-numerai/numerai"
	"gotest.tools/assert"
	"testing"
	"time"
)

func Test_Rounds(t *testing.T) {
	curr, err := numerai.QueryCurrentRound()
	assert.NilError(t, err)
	now := time.Now()
	assert.Assert(t, curr.OpenTime.Before(now))
	assert.Assert(t, curr.CloseTime.After(now))
	prev, err := numerai.QueryRoundInfo(curr.Number - 1)
	assert.NilError(t, err)
	assert.Assert(t, prev.Number+1 == curr.Number)
	assert.Assert(t, prev.OpenTime.Before(now))
	assert.Assert(t, prev.CloseTime.Before(now))
}
