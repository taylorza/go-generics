package rbtree_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/taylorza/go-generics/pkg/container/rbtree"
)

func Test_RBInsertAndSearch(t *testing.T) {
	tree := rbtree.New[int, string]()
	tree.Add(0, "Zero")
	tree.Add(1, "One")
	tree.Add(2, "Two")
	tree.Add(3, "Three")

	cases := []struct {
		desc          string
		key           int
		expectedOk    bool
		expectedValue string
	}{
		{"TestZeroExists", 0, true, "Zero"},
		{"TestOneExists", 1, true, "One"},
		{"TestTwoExists", 2, true, "Two"},
		{"TestThreeExists", 3, true, "Three"},
		{"TestNineDoesNotExists", 9, false, ""},
	}

	for _, tc := range cases {
		v, ok := tree.Search(tc.key)
		if ok != tc.expectedOk || v != tc.expectedValue {
			t.Fatalf("%s expected (value, ok)=(%v, %v) got (%v, %v) for key %v", tc.desc, tc.expectedValue, tc.expectedOk, v, ok, tc.key)
		}
	}
}

func Test_RBRemove(t *testing.T) {
	tree := rbtree.New[int, string]()
	tree.Add(0, "Zero")
	tree.Add(1, "One")
	tree.Add(2, "Two")
	tree.Add(3, "Three")

	tree.Remove(2)

	cases := []struct {
		desc          string
		key           int
		expectedOk    bool
		expectedValue string
	}{
		{"TestZeroExists", 0, true, "Zero"},
		{"TestOneExists", 1, true, "One"},
		{"TestTwoExists", 2, false, ""},
		{"TestThreeExists", 3, true, "Three"},
	}

	for _, tc := range cases {
		v, ok := tree.Search(tc.key)
		if ok != tc.expectedOk || v != tc.expectedValue {
			t.Fatalf("%s expected (value, ok)=(%v, %v) got (%v, %v) for key %v", tc.desc, tc.expectedValue, tc.expectedOk, v, ok, tc.key)
		}
	}
}

func Test_RBIter(t *testing.T) {
	tree := rbtree.New[int, string]()
	tree.Add(2, "Two")
	tree.Add(0, "Zero")
	tree.Add(3, "Three")
	tree.Add(1, "One")

	expected := []string{"Zero", "One", "Two", "Three"}
	got := []string{}

	it := tree.Iter()
	for it.Next() {
		got = append(got, it.Value())
	}

	if len(expected) != len(got) {
		t.Fatalf("Expected %v got '%v'", expected, got)
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != got[i] {
			t.Fatalf("Expected %v got '%v'", expected, got)
		}
	}
}

func Test_RBIterChan(t *testing.T) {
	tree := rbtree.New[int, string]()
	tree.Add(2, "Two")
	tree.Add(0, "Zero")
	tree.Add(3, "Three")
	tree.Add(1, "One")

	expected := []string{"Zero", "One", "Two", "Three"}
	got := []string{}

	for item := range tree.IterChan() {
		got = append(got, item.Value())
	}

	if len(expected) != len(got) {
		t.Fatalf("Expected %v got '%v'", expected, got)
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != got[i] {
			t.Fatalf("Expected %v got '%v'", expected, got)
		}
	}
}

func Test_RBRandomInsert(t *testing.T) {
	rand.Seed(1)

	expected := sort.IntSlice{}

	tree := rbtree.New[int, int]()
	for i := 0; i < 100; i++ {
		k := rand.Intn(1000)
		if _, ok := tree.Search(k); !ok {
			tree.Add(k, i)
			expected = append(expected, k)
		}
	}
	sort.Sort(expected)

	got := []int{}
	it := tree.Iter()
	for it.Next() {
		got = append(got, it.Key())
	}

	if len(expected) != len(got) {
		t.Fatalf("Expected %v got '%v'", expected, got)
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != got[i] {
			t.Fatalf("Expected %v got '%v'", expected, got)
		}
	}
}

func Test_RBRandomRemove(t *testing.T) {
	rand.Seed(1)

	expected := sort.IntSlice{}

	tree := rbtree.New[int, int]()
	for i := 0; i < 100; i++ {
		k := rand.Intn(1000)
		if _, ok := tree.Search(k); !ok {
			tree.Add(k, i)
			expected = append(expected, k)
		}
	}
	sort.Sort(expected)

	// Randomly remove items from the tree and the expected list
	for i := 0; i < len(expected); {
		if rand.Float64() >= 0.5 {
			tree.Remove(expected[i])
			expected = append(expected[0:i], expected[i+1:]...)
		} else {
			i++
		}
	}

	got := []int{}
	it := tree.Iter()
	for it.Next() {
		got = append(got, it.Key())
	}

	if len(expected) != len(got) {
		t.Fatalf("Expected %v got '%v'", expected, got)
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != got[i] {
			t.Fatalf("Expected %v got '%v'", expected, got)
		}
	}
}
