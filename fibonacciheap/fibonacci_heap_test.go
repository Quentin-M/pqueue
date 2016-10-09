package fibonacciheap

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	pqueue := New()

	e, k := pqueue.Peek()
	assert.Nil(t, e)
	assert.EqualValues(t, math.Inf(-1), k)

	e, k = pqueue.Pop()
	assert.Nil(t, e)
	assert.EqualValues(t, math.Inf(-1), k)

	assert.Zero(t, pqueue.Length())
}

func TestPushPopPeek(t *testing.T) {
	pqueue := New()

	pqueue.Push("O", math.Inf(0))
	pqueue.Push("F", -2)
	pqueue.Push("O", 0)

	e, k := pqueue.Peek()
	assert.EqualValues(t, "F", e)
	assert.EqualValues(t, -2, k)
	e, k = pqueue.Pop()
	assert.EqualValues(t, "F", e)
	assert.EqualValues(t, -2, k)
	e, k = pqueue.Pop()
	assert.EqualValues(t, "O", e)
	assert.EqualValues(t, 0, k)
	e, k = pqueue.Pop()
	assert.EqualValues(t, "O", e)
	assert.True(t, math.IsInf(k, 1))

	assert.Zero(t, pqueue.Length())
}

func TestDecreaseKey(t *testing.T) {
	pqueue := New()

	pqueue.Push("B", 0)
	ie := pqueue.Push("A", 3)
	pqueue.DecreaseKey(ie, 1)
	pqueue.Push("R", 2)

	e, k := pqueue.Pop()
	assert.EqualValues(t, "B", e)
	assert.EqualValues(t, 0, k)

	e, k = pqueue.Pop()
	assert.EqualValues(t, "A", e)
	assert.EqualValues(t, 1, k)
}

func TestDecreaseKeyNonExisting(t *testing.T) {
	pqueue := New()

	ie := pqueue.Push("B", 0)
	pqueue.Pop()

	assert.Panics(t, func() {
		pqueue.DecreaseKey(ie, -1)
	})
}

func TestIncreaseKey(t *testing.T) {
	pqueue := New()

	ie := pqueue.Push("F", 0)
	assert.Panics(t, func() {
		pqueue.DecreaseKey(ie, 1)
	})

	e, k := pqueue.Pop()
	assert.EqualValues(t, e, "F")
	assert.EqualValues(t, k, 0)
}

func TestHasGet(t *testing.T) {
	pqueue := New()

	e, k := pqueue.Get(nil)
	assert.EqualValues(t, e, nil)
	assert.True(t, math.IsInf(k, -1))

	ie := pqueue.Push("F", 0)
	assert.True(t, pqueue.Has(ie))

	e, k = pqueue.Get(ie)
	assert.EqualValues(t, e, "F")
	assert.EqualValues(t, k, 0)

	pqueue.Pop()
	assert.False(t, pqueue.Has(ie))
}

func TestDelete(t *testing.T) {
	pqueue := New()

	pqueue.Push("O", 5)
	ie := pqueue.Push("F", -2)
	pqueue.Delete(ie)

	e, k := pqueue.Pop()
	assert.EqualValues(t, "O", e)
	assert.EqualValues(t, 5, k)
}

func TestDeleteNonExisting(t *testing.T) {
	pqueue := New()

	assert.Panics(t, func() {
		pqueue.Delete(nil)
	})

	ie := pqueue.Push("O", 5)
	pqueue.Pop()

	assert.Panics(t, func() {
		pqueue.Delete(ie)
	})
}

func TestClear(t *testing.T) {
	pqueue := New()

	pqueue.Push("O", 5)
	pqueue.Clear()

	assert.Zero(t, pqueue.Length())

	e, k := pqueue.Pop()
	assert.Nil(t, e)
	assert.True(t, math.IsInf(k, -1))
}

func TestMinusInfKey(t *testing.T) {
	pqueue := New()

	assert.Panics(t, func() {
		pqueue.Push("F", math.Inf(-1))
	})
	assert.Zero(t, pqueue.Length())

	ie := pqueue.Push("F", 0)
	assert.Panics(t, func() {
		pqueue.DecreaseKey(ie, math.Inf(-1))
	})

	e, k := pqueue.Pop()
	assert.EqualValues(t, e, "F")
	assert.EqualValues(t, k, 0)
}

type testItem struct {
	node  interface{}
	key   float64
	value float64
}

type ByKey []testItem

func (s ByKey) Len() int {
	return len(s)
}

func (s ByKey) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByKey) Less(i, j int) bool {
	return s[i].key < s[j].key
}

func TestPushPopDecreaseBig(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var sortedItems []testItem
	pqueue := New()

	for i := 0; i < 1000; i++ {
		modeSelector := rand.Float64()

		if len(sortedItems) == 0 || modeSelector < 0.50 {
			// Push (50%).
			randomValue := rand.Float64() * 1000
			item := testItem{key: randomValue, value: randomValue}
			item.node = pqueue.Push(randomValue, randomValue)
			sortedItems = append(sortedItems, item)
			sort.Sort(ByKey(sortedItems))
		} else if modeSelector < 0.85 {
			e, k := pqueue.Pop()
			if e.(float64) != sortedItems[0].value || k != sortedItems[0].key {
				fmt.Printf("k: %.2f | v: %.2f - expected k: %.2f | v: %.2f\n", k, e.(float64), sortedItems[0].key, sortedItems[0].value)
				assert.FailNow(t, "Invalid Pop()")
			}
			sortedItems = sortedItems[1:]
		} else {
			// DecreaseKey (15%).
			index := rand.Intn(len(sortedItems))
			randomValue := rand.Float64() * sortedItems[index].key
			pqueue.DecreaseKey(sortedItems[index].node, randomValue)
			sortedItems[index].key = randomValue
			sort.Sort(ByKey(sortedItems))
		}
	}

	for _, item := range sortedItems {
		e, k := pqueue.Pop()
		if e.(float64) != item.value || k != item.key {
			assert.FailNow(t, "Invalid Pop()")
		}
	}
}
