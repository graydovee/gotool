package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestSwap(t *testing.T) {
	arr := []int{2, 3}
	Swap(arr, 0, 1)
	e := arr[0] == 3 && arr[1] == 2
	if !e {
		t.Errorf("swap error")
	}

	slice := Slice[int]{1, 2}
	slice.Swap(0, 1)
	e = slice[0] == 2 && slice[1] == 1
	if !e {
		t.Errorf("swap error")
	}
}

func TestRange(t *testing.T) {
	slice := make(Slice[int], 10000)
	for i := 0; i < 10000; i++ {
		slice[i] = i
	}

	now := time.Now()
	j := 0
	for _, i := range slice {
		j += i
	}
	c1 := time.Since(now)
	fmt.Println(c1)

	now = time.Now()
	j = 0
	for i := range slice.Iter() {
		j += i
	}
	c2 := time.Since(now)
	fmt.Println(c2)

	now = time.Now()
	j = 0
	iterator := slice.Iterator()
	for i := range iterator.C {
		j += i
	}
	c3 := time.Since(now)
	fmt.Println(c3)

	now = time.Now()
	j = 0
	slice.Foreach(func(i int, e int) bool {
		j += e
		return true
	})
	c4 := time.Since(now)
	fmt.Println(c4)

	fmt.Println(c2 / c1)
}
