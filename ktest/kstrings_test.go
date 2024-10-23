package ktest

import (
	"testing"

	"github.com/khan-lau/kutils/container/kobjs"
	"github.com/khan-lau/kutils/container/kstrings"
)

func TestKStrings(t *testing.T) {
	formatTuple1 := kstrings.SliceFormat("aaa %s bbb %d cc {}", "a", 1, "c")
	// formatTuple2 := kstrings.SliceFormat("aaa {} bbb {} ccc {}", "a", 1, "c")
	t.Log(kobjs.ObjectToJson5WithoutFunc(formatTuple1))
	// t.Log(kobjs.ObjectToJson5WithoutFunc(formatTuple2))

	// t.Log(kstrings.FormatString("aaa %s bbb %d", "a", 1))
	// t.Log(kstrings.FormatString("aaa {} bbb {} ccc {}", "a", 1, "c"))

	t.Log(kstrings.FormatString("aaa %s bbb %d ccc {}", "a", 1, "c"))
}
