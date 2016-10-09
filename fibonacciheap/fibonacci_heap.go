// Package fibonacciheap contains a Fibonacci implementation of priority queues.
package fibonacciheap

import (
	"container/list"
	"math"

	"github.com/quentin-m/pqueue"
)

type fibonacciHeap struct {
	roots  *list.List
	min    *node
	length int
}

type node struct {
	key   float64
	value interface{}

	self     *list.Element
	parent   *node
	children *list.List

	isQueued bool
	marked   bool
}

func toNode(element *list.Element) *node {
	if element == nil {
		return nil
	}
	return element.Value.(*node)
}

// New initializes a new Fibonacci heap.
//
// The amortized running time of most of these methods is O(1), making it a very
// fast data structure. Several have an actual running time of O(1).
// Pop() and Delete() have O(log n) amortized running times because they do
// the heap consolidation.
//
// Note that this implementation is not synchronized. If multiple threads
// access a set concurrently, and at least one of the threads modifies the set,
// it must be synchronized externally. This is typically accomplished by
// synchronizing on some object that naturally encapsulates the set.
func New() pqueue.PriorityQueue {
	return &fibonacciHeap{
		roots: list.New(),
	}
}

// Push inserts a new prioritized value into the heap.
//
// The priority must be in range (-inf,+inf], otherwise, a panic occurs.
// The returned value can be specified as a parameter to other methods such as
// Has(), Get(), DecreaseKey() or Delete().
//
// Time complexity: O(1) actual.
func (fh *fibonacciHeap) Push(value interface{}, key float64) interface{} {
	// Verify precondition.
	if math.IsInf(key, -1) {
		panic("key must be in range (-inf,+inf]")
	}

	node := &node{
		key:      key,
		value:    value,
		children: list.New(),
		isQueued: true,
	}

	// Add new node to the root and node lists.
	node.self = fh.roots.PushBack(node)
	fh.length++

	// Set the node as the minimum one if necessary.
	if fh.min == nil || fh.min.key > node.key {
		fh.min = node
	}

	return node
}

// Peek returns the value of maximum priority (i.e. minimum value) without
// removing it from the queue.
// Time complexity: O(1) actual.
func (fh *fibonacciHeap) Peek() (interface{}, float64) {
	if fh.min == nil {
		return nil, math.Inf(-1)
	}
	return fh.min.value, fh.min.key
}

// Pop returns the value of maximum priority (i.e. minimum value) and removes it
// from the queue.
// Time complexity: O(log n) amortized.
func (fh *fibonacciHeap) Pop() (interface{}, float64) {
	if fh.min == nil {
		return nil, math.Inf(-1)
	}

	// Unparent each child of the minimum node and add them to the root list.
	for e := fh.min.children.Front(); e != nil; e = e.Next() {
		child := toNode(e)
		child.self = fh.roots.PushBack(child)
		child.parent = nil
		child.marked = false
	}
	fh.min.children.Init()

	// Delete the minimum node from the root list and decrease the heap's length.
	fh.roots.Remove(fh.min.self)
	fh.length--

	// Store the current minimum node for later return and mark it as non-queued.
	pmin := fh.min
	pmin.isQueued = false

	// Consolidate the heap and reset the minimum.
	fh.consolidate()
	fh.resetMin()

	return pmin.value, pmin.key
}

// Has returns whether the given element is present in the queue.
// Time complexity: O(1) actual.
func (fh *fibonacciHeap) Has(iNode interface{}) bool {
	if iNode == nil || !iNode.(*node).isQueued {
		return false
	}
	return true
}

// Get returns the value and the priority of the given element.
// Time complexity: O(1) actual.
func (fh *fibonacciHeap) Get(iNode interface{}) (interface{}, float64) {
	if iNode == nil {
		return nil, math.Inf(-1)
	}
	node := iNode.(*node)
	return node.value, node.key
}

// DecreaseKey increases the priority of the given element.
//
// The element must be present in the queue, otherwise, a panic occurs.
// The priority must be in range (-inf,old priority], otherwise, a panic occurs.
// The returned value can be specified as a parameter to other methods such as
// Has(), Get(), DecreaseKey() or Delete().
//
// Time complexity: O(1) amortized.
func (fh *fibonacciHeap) DecreaseKey(iNode interface{}, key float64) {
	if math.IsInf(key, -1) {
		panic("key must be in range (-inf,+inf]")
	}
	fh.decreaseKey(iNode, key)
}

// Internal version of DecreaseKey(). The public version calls it after
// verifying that the key isn't -inf, which is reserved for Delete().
func (fh *fibonacciHeap) decreaseKey(iNode interface{}, key float64) {
	if iNode == nil {
		panic("nil value given")
	}
	node := iNode.(*node)

	// Verify precondition.
	if key > node.key {
		panic("new key is greater than current key")
	}
	if !node.isQueued {
		panic("node isn't queued")
	}

	// Assign new key.
	node.key = key

	// Cut the node from its parent and cascade the cuts upward if the heap
	// property just became violated.
	if node.parent != nil && node.key < node.parent.key {
		p := node.parent
		fh.cut(node)
		fh.cascadingCut(p)
	}

	// Set the minimum of the heap if necessary.
	if node.key < fh.min.key {
		fh.min = node
	}
}

// Delete removes the given element from the queue.
//
// The element must be present in the queue, otherwise, a panic occurs.
//
// Time complexity: O(log n) amortized.
func (fh *fibonacciHeap) Delete(node interface{}) {
	fh.decreaseKey(node, math.Inf(-1))
	fh.Pop()
}

// Length returns the number of elements in the queue.
// Time complexity: O(1)
func (fh *fibonacciHeap) Length() int {
	return fh.length
}

// Clear resets the list to its initial state, forgetting all its data.
// Time complexity: O(1)
func (fh *fibonacciHeap) Clear() {
	fh.roots.Init()
	fh.min = nil
	fh.length = 0
}

func (fh *fibonacciHeap) maxDegree() int {
	phi := (1.0 + math.Sqrt(5.0)) / 2.0
	return (int)(math.Ceil(math.Log((float64)(fh.length)) / math.Log(phi)))
}

func (fh *fibonacciHeap) consolidate() {
	if fh.length == 0 {
		return
	}

	treeByDegree := make(map[int]*list.Element, fh.maxDegree())
	// For each root element.
	for treeElement := fh.roots.Front(); treeElement != nil; {
		treeNode := toNode(treeElement)
		treeDegree := treeNode.children.Len()

		// If there isn't any other tree of that degree, assign it and continue.
		otherTreeElement, otherTreeExists := treeByDegree[treeDegree]
		if !otherTreeExists || otherTreeElement == treeElement {
			treeByDegree[treeDegree] = treeElement
			treeElement = treeElement.Next()
			continue
		}

		// There is another tree with that degree, link the one that has a higher
		// key into the other one, increase the degree of the tree and redo if
		// there is another tree of the new degree.
		for ; otherTreeElement != nil; otherTreeElement = treeByDegree[treeDegree] {
			otherTreeNode := toNode(otherTreeElement)
			if treeNode.key <= otherTreeNode.key {
				fh.link(otherTreeNode, treeNode)
			} else {
				fh.link(treeNode, otherTreeNode)
				treeElement = otherTreeElement
				treeNode = otherTreeNode
			}

			treeByDegree[treeDegree] = nil
			treeDegree++
		}
		treeByDegree[treeDegree] = treeElement
	}
}

func (fh *fibonacciHeap) resetMin() {
	min := fh.roots.Front()

	for tree := fh.roots.Front(); tree != nil; tree = tree.Next() {
		if toNode(tree).key < toNode(min).key {
			min = tree
		}
	}

	fh.min = toNode(min)
}

func (fh *fibonacciHeap) link(node, newParent *node) {
	fh.roots.Remove(node.self)
	node.parent = newParent
	node.marked = false
	node.self = newParent.children.PushBack(node)
}

func (fh *fibonacciHeap) cut(node *node) {
	node.parent.children.Remove(node.self)
	node.parent = nil
	node.marked = false
	node.self = fh.roots.PushBack(node)
}

func (fh *fibonacciHeap) cascadingCut(node *node) {
	parent := node.parent
	if parent == nil {
		return
	}

	if node.marked {
		fh.cut(node)
		fh.cascadingCut(parent)
	} else {
		node.marked = true
	}
}
