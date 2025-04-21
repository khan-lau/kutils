package ktest

import (
	"testing"

	"github.com/khan-lau/kutils/container/kring"
)

func Test_Ring_Len(t *testing.T) {
	// Create a new ring of size 4
	r := kring.New[int](4)

	// Print out its length
	t.Logf("len: %d", r.Count())

	// Output:
	// 4
	if r.Count() != 4 {
		t.Errorf("expected len: %d, got: %d", 4, r.Count())
	}
}

func Test_Ring_Next(t *testing.T) {
	// Create a new ring of size 5
	r := kring.New[int](5)

	// Get the length of the ring
	n := r.Count()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Value = i
		r = r.Next()
	}

	// Iterate through the ring and print its contents
	for j := 0; j < n; j++ {
		t.Logf("val: %d", r.Value)
		r = r.Next()
	}

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func Test_Ring_Prev(t *testing.T) {
	// Create a new ring of size 5
	r := kring.New[int](5)

	// Get the length of the ring
	n := r.Count()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Value = i
		r = r.Next()
	}

	// Iterate through the ring backwards and print its contents
	for j := 0; j < n; j++ {
		r = r.Prev()
		t.Logf("val: %d", r.Value)
	}

	// Output:
	// 4
	// 3
	// 2
	// 1
	// 0
}

func Test_Ring_Do(t *testing.T) {
	// Create a new ring of size 5
	r := kring.New[int](5)

	// Get the length of the ring
	n := r.Count()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Value = i
		r = r.Next()
	}

	// Iterate through the ring and print its contents
	r.Do(func(p any) {
		t.Logf("val: %d", p.(int))
	})

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func Test_Ring_Move(t *testing.T) {
	// Create a new ring of size 5
	r := kring.New[int](5)

	// Get the length of the ring
	n := r.Count()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Value = i
		r = r.Next()
	}

	// Move the pointer forward by three steps
	r = r.Move(3)

	// Iterate through the ring and print its contents
	r.Do(func(p any) {
		t.Logf("val: %d", p.(int))
	})

	// Output:
	// 3
	// 4
	// 0
	// 1
	// 2
}

func Test_Ring_Link(t *testing.T) {
	// Create two rings, r and s, of size 2
	r := kring.New[int](2)
	s := kring.New[int](2)

	// Get the length of the ring
	lr := r.Count()
	ls := s.Count()

	// Initialize r with 0s
	for i := 0; i < lr; i++ {
		r.Value = 0
		r = r.Next()
	}

	// Initialize s with 1s
	for j := 0; j < ls; j++ {
		s.Value = 1
		s = s.Next()
	}

	// Link ring r and ring s
	rs := r.Link(s)

	// Iterate through the combined ring and print its contents
	rs.Do(func(p any) {
		t.Logf("val: %d", p.(int))
	})

	// Output:
	// 0
	// 0
	// 1
	// 1
}

func Test_Ring_Unlink(t *testing.T) {
	// Create a new ring of size 6
	r := kring.New[int](6)

	// Get the length of the ring
	n := r.Count()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Value = i
		r = r.Next()
	}

	// Unlink three elements from r, starting from r.Next()
	r.Unlink(3)

	// Iterate through the remaining ring and print its contents
	r.Do(func(p any) {
		t.Logf("val: %d", p.(int))
	})

	// Output:
	// 0
	// 4
	// 5
}
