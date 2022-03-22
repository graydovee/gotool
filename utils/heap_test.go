package utils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMaxHeap(t *testing.T) {
	var h MaxHeap[int]
	num := 1000
	max := 10000
	rand.Seed(time.Now().UnixMicro())
	for i := 0; i < num; i++ {
		h.Add(rand.Int() % max)
	}

	for i := range h.Iter() {
		fmt.Print(i, ", ")
	}

	var pre = max
	fmt.Print("\n", pre, ", ")
	for i := 0; i < num; i++ {
		n := h.Pop()
		fmt.Print(n, ", ")
		if n > pre {
			t.Errorf("heap error")
			return
		}
		pre = n
	}
}

func TestMinHeap(t *testing.T) {
	var h MinHeap[int]
	num := 1000
	rand.Seed(time.Now().UnixMicro())
	for i := 0; i < num; i++ {
		h.Add(rand.Int() % 10000)
	}

	for i := range h.Iter() {
		fmt.Print(i, ", ")
	}
	fmt.Print("\n")

	var pre = 0
	for i := 0; i < num; i++ {
		n := h.Pop()
		if n < pre {
			t.Errorf("heap error")
			return
		}
		fmt.Print(pre, ", ")
		pre = n
	}
}
