package tests

import (
	"fmt"
	"github.com/sudachen/go-numerai/numerai"
	"testing"
)

func Test_Leaderboard(t *testing.T) {
	q := numerai.Leaderboard()
	fmt.Println(q.Head(10))
}
