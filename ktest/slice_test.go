package ktest

import (
	"testing"

	"github.com/khan-lau/kutils/container/klists"
	"github.com/khan-lau/kutils/container/kslices"
)

func Test_SplitSlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	limit := 16
	for i := 0; i < limit; i++ {
		result := kslices.SplitSliceByLimit(s, i)
		t.Log(result) // 输出：[[1 2] [3 4] [5 6] [7 8] [9 10]]
	}
}

func Test_SplitKList(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	kl := klists.New[int]()
	kl.PushBackSlice(s...)

	limit := 16
	for i := 0; i < limit; i++ {
		result := klists.SplitKList(kl, i)
		t.Log(result) // 输出：[[1 2 3 4] [5 6 7 8] [9 10]]
	}

}
