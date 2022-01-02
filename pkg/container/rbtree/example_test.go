package rbtree_test

import (
	"fmt"

	"github.com/taylorza/go-generics/pkg/container/rbtree"
)

func ExampleTree_Add() {
	tree := rbtree.New[int, string]()
	tree.Add(0, "Zero")
	tree.Add(1, "One")
	tree.Add(2, "Two")
	tree.Add(3, "Three")
	tree.Add(4, "Four")
	tree.Add(5, "Five")
	tree.Add(6, "Six")
	tree.Add(7, "Seven")
	tree.Add(8, "Eight")
	tree.Add(9, "Nine")

	for i := 0; i < 10; i++ {
		if s, ok := tree.Search(i); ok {
			fmt.Printf("%v - %v\n", i, s)
		}
	}
	// Output:
	// 0 - Zero
	// 1 - One
	// 2 - Two
	// 3 - Three
	// 4 - Four
	// 5 - Five
	// 6 - Six
	// 7 - Seven
	// 8 - Eight
	// 9 - Nine
}

func ExampleTree_Iter() {
	tree := rbtree.New[int, string]()
	tree.Add(0, "Zero")
	tree.Add(1, "One")
	tree.Add(2, "Two")
	tree.Add(3, "Three")
	tree.Add(4, "Four")
	tree.Add(5, "Five")
	tree.Add(6, "Six")
	tree.Add(7, "Seven")
	tree.Add(8, "Eight")
	tree.Add(9, "Nine")

	it := tree.Iter()
	for it.Next() {
		fmt.Printf("%v - %v\n", it.Key(), it.Value())
	}
	// Output:
	// 0 - Zero
	// 1 - One
	// 2 - Two
	// 3 - Three
	// 4 - Four
	// 5 - Five
	// 6 - Six
	// 7 - Seven
	// 8 - Eight
	// 9 - Nine
}
