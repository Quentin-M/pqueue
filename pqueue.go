// Package pqueue contains an interface defining priority queues.
package pqueue

// PriorityQueue is an interface that defines several common operations
// of the abstract data structure of the same name.
type PriorityQueue interface {
	Push(interface{}, float64) interface{}
	Peek() (interface{}, float64)
	Pop() (interface{}, float64)
	Has(interface{}) bool
	Get(interface{}) (interface{}, float64)
	DecreaseKey(interface{}, float64)
	Delete(interface{})
	Length() int
	Clear()
}
