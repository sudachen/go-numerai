package tests

import (
	"fmt"
	"github.com/sudachen/go-numerai/nuai"
	"testing"
)

func Test_Leaderboard(t *testing.T) {
	n := nuai.Nuai{}
	q := n.GetLeaderboard(99)
	fmt.Println(q.Slice(33,53).Display())
}
