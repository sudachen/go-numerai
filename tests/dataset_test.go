package tests

import (
	"fmt"
	"github.com/sudachen/go-numerai/numerai"
	"gotest.tools/assert"
	"testing"
)

func Test_DownloadDataset(t *testing.T) {
	q, err := numerai.Training.First(100).Collect()
	assert.NilError(t, err)
	fmt.Println(q.Len(), len(q.Names()))
	fmt.Println(q.Head(5))
}

func Test_DownloadDataset1(t *testing.T) {
	q, err := numerai.Tournament.Count()
	assert.NilError(t, err)
	fmt.Println(q)
}

func Test_TraningDataset(t *testing.T) {
	fmt.Println(numerai.Training.First(5).LuckyCollect())
}
