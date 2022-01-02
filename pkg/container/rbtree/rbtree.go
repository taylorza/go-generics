package rbtree

import (
	"constraints"
)

type color byte

const (
	black color = iota
	red
)

// Node interface represents an entry in the tree.
type Node[K constraints.Ordered, V any] interface {
	// Key returns the key for the node in the tree.
	Key() K
	// Value returns the values associated with the node in the tree.
	Value() V
}

type node[K constraints.Ordered, V any] struct {
	k K
	v V
	p *node[K, V]
	l *node[K, V]
	r *node[K, V]
	c color
}

// Key returns the key for the node in the tree.
func (n *node[K, V]) Key() K {
	return n.k
}

// Value returns the values associated with the node in the tree.
func (n *node[K, V]) Value() V {
	return n.v
}

func (n *node[K, V]) grandParent() *node[K, V] {
	return n.p.p
}

func (n *node[K, V]) sibling() *node[K, V] {
	if n == n.p.l {
		return n.p.r
	} else {
		return n.p.l
	}
}

func (n *node[K, V]) uncle() *node[K, V] {
	return n.p.sibling()
}

func nodeColor[K constraints.Ordered, V any](n *node[K, V]) color {
	if n == nil {
		return black
	}
	return n.c
}

func maxNode[K constraints.Ordered, V any](n *node[K, V]) *node[K, V] {
	for n.r != nil {
		n = n.r
	}
	return n
}

func newNode[K constraints.Ordered, V any](k K, v V, c color, l, r *node[K, V]) *node[K, V] {
	n := &node[K, V]{k: k, v: v, c: c, l: l, r: r}
	if l != nil {
		l.p = n
	}
	if r != nil {
		r.p = n
	}
	return n
}

// Tree is a red-black tree which is a self-balancing binary search tree.
type Tree[K constraints.Ordered, V any] struct {
	root  *node[K, V]
	count int
}

// New returns a new instance of a tree.
func New[K constraints.Ordered, V any]() *Tree[K, V] {
	return &Tree[K, V]{}
}

// Add adds the key/value pair to the tree.
// If the key already exists in the tree, the value associated with the key is updated to the new value.
func (t *Tree[K, V]) Add(key K, value V) {
	nn := newNode(key, value, red, nil, nil)
	if t.root == nil {
		t.root = nn
	} else {
		n := t.root
		for {
			if key == n.k {
				n.v = value
				return
			} else if key < n.k {
				if n.l == nil {
					n.l = nn
					break
				} else {
					n = n.l
				}
			} else {
				if n.r == nil {
					n.r = nn
					break
				} else {
					n = n.r
				}
			}
		}
		nn.p = n
	}
	t.insertCase1(nn)
	t.count++
}

// Remove deletes the item from the tree.
// If the item does not exist the function returns false
func (t *Tree[K, V]) Remove(key K) bool {
	n := t.lookup(key)
	if n == nil {
		return false
	}

	if n.l != nil && n.r != nil {
		// Copy key/value from predecessor and then delete it instead
		pred := maxNode(n.l)
		n.k = pred.k
		n.v = pred.v
		n = pred
	}

	child := n.l
	if n.r != nil {
		child = n.r
	}
	if nodeColor(n) == black {
		n.c = nodeColor(child)
		t.deleteCase1(n)
	}
	t.replaceNode(n, child)
	if n.p == nil && child != nil {
		child.c = black
	}
	t.count--
	return true
}

// Search searches the tree for an item with the specified key and returns the value if it exists.
// ok is false if the key does not exist.
func (t *Tree[K, V]) Search(key K) (value V, ok bool) {
	n := t.lookup(key)
	if n == nil {
		var z V
		return z, false
	}
	return n.v, true
}

// Len returns the number of items in the tree.
func (t *Tree[K, V]) Len() int {
	return t.count
}

func (t *Tree[K, V]) lookup(key K) *node[K, V] {
	n := t.root
	for n != nil {
		if key == n.k {
			return n
		} else if key < n.k {
			n = n.l
		} else {
			n = n.r
		}
	}
	return nil
}

func (t *Tree[K, V]) rotateLeft(n *node[K, V]) {
	r := n.r
	t.replaceNode(n, r)
	n.r = r.l
	if r.l != nil {
		r.l.p = n
	}
	r.l = n
	n.p = r
}

func (t *Tree[K, V]) rotateRight(n *node[K, V]) {
	l := n.l
	t.replaceNode(n, l)
	n.l = l.r
	if l.r != nil {
		l.r.p = n
	}
	l.r = n
	n.p = l
}

func (t *Tree[K, V]) replaceNode(n, nn *node[K, V]) {
	if n.p == nil {
		t.root = nn
	} else {
		if n == n.p.l {
			n.p.l = nn
		} else {
			n.p.r = nn
		}
	}
	if nn != nil {
		nn.p = n.p
	}
}

func (t *Tree[K, V]) insertCase1(n *node[K, V]) {
	if n.p == nil {
		n.c = black
	} else {
		t.insertCase2(n)
	}
}

func (t *Tree[K, V]) insertCase2(n *node[K, V]) {
	if nodeColor(n.p) == black {
		return
	} else {
		t.insertCase3(n)
	}
}

func (t *Tree[K, V]) insertCase3(n *node[K, V]) {
	if nodeColor(n.uncle()) == red {
		n.p.c = black
		n.uncle().c = black
		n.grandParent().c = red
		t.insertCase1(n.grandParent())
	} else {
		t.insertCase4(n)
	}
}

func (t *Tree[K, V]) insertCase4(n *node[K, V]) {
	if n == n.p.r && n.p == n.grandParent().l {
		t.rotateLeft(n.p)
		n = n.l
	} else if n == n.p.l && n.p == n.grandParent().r {
		t.rotateRight(n.p)
		n = n.r
	}
	t.insertCase5(n)
}

func (t *Tree[K, V]) insertCase5(n *node[K, V]) {
	n.p.c = black
	n.grandParent().c = red
	if n == n.p.l && n.p == n.grandParent().l {
		t.rotateRight(n.grandParent())
	} else {
		t.rotateLeft(n.grandParent())
	}
}

func (t *Tree[K, V]) deleteCase1(n *node[K, V]) {
	if n.p == nil {
		return
	} else {
		t.deleteCase2(n)
	}
}

func (t *Tree[K, V]) deleteCase2(n *node[K, V]) {
	if nodeColor(n.sibling()) == red {
		n.p.c = red
		n.sibling().c = black
		if n == n.p.l {
			t.rotateLeft(n.p)
		} else {
			t.rotateRight(n.p)
		}
	}
	t.deleteCase3(n)
}

func (t *Tree[K, V]) deleteCase3(n *node[K, V]) {
	if nodeColor(n.p) == black &&
		nodeColor(n.sibling()) == black &&
		nodeColor(n.sibling().l) == black &&
		nodeColor(n.sibling().r) == black {
		n.sibling().c = red
		t.deleteCase1(n.p)
	} else {
		t.deleteCase4(n)
	}
}

func (t *Tree[K, V]) deleteCase4(n *node[K, V]) {
	if nodeColor(n.p) == red &&
		nodeColor(n.sibling()) == black &&
		nodeColor(n.sibling().l) == black &&
		nodeColor(n.sibling().r) == black {
		n.sibling().c = red
		n.p.c = black
	} else {
		t.deleteCase5(n)
	}
}

func (t *Tree[K, V]) deleteCase5(n *node[K, V]) {
	if n == n.p.l &&
		nodeColor(n.sibling()) == black &&
		nodeColor(n.sibling().l) == red &&
		nodeColor(n.sibling().r) == black {
		n.sibling().c = red
		n.sibling().l.c = black
		t.rotateRight(n.sibling())
	} else if n == n.p.r &&
		nodeColor(n.sibling()) == black &&
		nodeColor(n.sibling().r) == red &&
		nodeColor(n.sibling().l) == black {
		n.sibling().c = red
		n.sibling().r.c = black
		t.rotateLeft(n.sibling())
	}
	t.deleteCase6(n)
}

func (t *Tree[K, V]) deleteCase6(n *node[K, V]) {
	n.sibling().c = nodeColor(n.p)
	n.p.c = black
	if n == n.p.l {
		n.sibling().r.c = black
		t.rotateLeft(n.p)
	} else {
		n.sibling().l.c = black
		t.rotateRight(n.p)
	}
}

// IterChan returns a channel used to iterate through all the items in the tree.
func (t *Tree[K, V]) IterChan() <-chan Node[K, V] {
	if t.root == nil {
		return nil
	}
	ch := make(chan Node[K, V])
	go func() {
		var s []*node[K, V]
		curr := t.root
		for len(s) > 0 || curr != nil {
			if curr != nil {
				s = append(s, curr)
				curr = curr.l
			} else {
				n := s[len(s)-1]
				s = s[:len(s)-1]
				ch <- n
				curr = n.r
			}
		}
		close(ch)
	}()
	return ch
}

// Iterator represents an iterator used to iterate items in the tree.
type Iterator[K constraints.Ordered, V any] struct {
	r *node[K, V]
	c *node[K, V]
	n *node[K, V]
	s []*node[K, V]
}

// Iter returns a new instance of an Iteraator that can be used to iterate through the items in the tree.
func (t *Tree[K, V]) Iter() *Iterator[K, V] {
	return &Iterator[K, V]{r: t.root, c: t.root}
}

// Reset resets the iterator to the begining.
func (it *Iterator[K, V]) Reset() {
	it.c = it.r
	it.n = nil
	it.s = it.s[:0]
}

// Next moves the iterator forward, returning false when there are no more items to iterate.
// Next must be called before trying to access the key or value through the iterator
func (it *Iterator[K, V]) Next() bool {
	if it.r == nil {
		return false
	}
	for len(it.s) > 0 || it.c != nil {
		if it.c != nil {
			it.s = append(it.s, it.c)
			it.c = it.c.l
		} else {
			it.n = it.s[len(it.s)-1]
			it.s = it.s[:len(it.s)-1]
			it.c = it.n.r
			return true
		}
	}
	return false
}

// Key returns the key of the current item in the iteration.
func (it *Iterator[K, V]) Key() K {
	if it.n == nil {
		panic("Key called without calling Next")
	}
	return it.n.k
}

// Value returns the value associated with the current item in the iteration.
func (it *Iterator[K, V]) Value() V {
	if it.n == nil {
		panic("Value called without calling Next")
	}
	return it.n.v
}
