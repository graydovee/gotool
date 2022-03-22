package utils

import (
	"golang.org/x/exp/constraints"
)

type Heap[E constraints.Ordered] interface {
	Iterable[E]
	Add(e E)
	Peek() E
	Pop() E
	Len() int
}

// MaxHeap minHeap
type MaxHeap[E constraints.Ordered] struct {
	nodes Slice[E]
}

func (h *MaxHeap[E]) Iterator() *Iterator[E] {
	return h.nodes.Iterator()
}

func (h *MaxHeap[E]) Iter() <-chan E {
	return h.nodes.Iter()
}

func (h *MaxHeap[E]) Foreach(f func(index int, element E) bool) {
	h.nodes.Foreach(f)
}

func (h *MaxHeap[E]) Len() int {
	return len(h.nodes)
}

func (h *MaxHeap[E]) Add(e E) {
	h.nodes = append(h.nodes, e)
	h.shiftUp(len(h.nodes) - 1)
}

func (h *MaxHeap[E]) Peek() E {
	return h.nodes[0]
}

func (h *MaxHeap[E]) Pop() E {
	v := h.nodes[0]
	h.nodes[0] = h.nodes[len(h.nodes)-1]
	h.nodes = h.nodes[:len(h.nodes)-1]
	h.shiftDown(0)
	return v
}

func (h *MaxHeap[E]) shiftDown(n int) {
	if n >= len(h.nodes) {
		return
	}

	cl := n << 1
	if cl >= len(h.nodes) {
		return
	}
	bigger := cl

	cr := n<<1 | 1
	if cr < len(h.nodes) && h.nodes[cl] < h.nodes[cr] {
		bigger = cr
	}

	if h.nodes[n] < h.nodes[bigger] {
		Swap(h.nodes, n, bigger)
		h.shiftDown(bigger)
	}
}

func (h *MaxHeap[E]) shiftUp(n int) {
	if n == 0 {
		return
	}
	p := n / 2
	if h.nodes[n] > h.nodes[p] {
		Swap(h.nodes, n, p)
		h.shiftUp(p)
	}
}

// MinHeap minHeap
type MinHeap[E constraints.Ordered] struct {
	nodes Slice[E]
}

func (h *MinHeap[E]) Iterator() *Iterator[E] {
	return h.nodes.Iterator()
}

func (h *MinHeap[E]) Iter() <-chan E {
	return h.nodes.Iter()
}

func (h *MinHeap[E]) Foreach(f func(index int, element E) bool) {
	h.nodes.Foreach(f)
}

func (h *MinHeap[E]) Len() int {
	return len(h.nodes)
}

func (h *MinHeap[E]) Add(e E) {
	h.nodes = append(h.nodes, e)
	h.shiftUp(len(h.nodes) - 1)
}

func (h *MinHeap[E]) Peek() E {
	return h.nodes[0]
}

func (h *MinHeap[E]) Pop() E {
	v := h.nodes[0]
	h.nodes[0] = h.nodes[len(h.nodes)-1]
	h.nodes = h.nodes[:len(h.nodes)-1]
	h.shiftDown(0)
	return v
}

func (h *MinHeap[E]) shiftDown(n int) {
	if n >= len(h.nodes) {
		return
	}

	cl := n << 1
	if cl >= len(h.nodes) {
		return
	}
	smaller := cl

	cr := n<<1 | 1
	if cr < len(h.nodes) && h.nodes[cl] > h.nodes[cr] {
		smaller = cr
	}

	if h.nodes[n] > h.nodes[smaller] {
		Swap(h.nodes, n, smaller)
		h.shiftDown(smaller)
	}
}

func (h *MinHeap[E]) shiftUp(n int) {
	if n == 0 {
		return
	}
	p := n / 2
	if h.nodes[n] < h.nodes[p] {
		Swap(h.nodes, n, p)
		h.shiftUp(p)
	}
}
