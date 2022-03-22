package utils

import (
	"testing"
)

func TestSwap(t *testing.T) {
	arr := []int{2, 3}
	Swap(arr, 0, 1)
	e := arr[0] == 3 && arr[1] == 2
	if !e {
		t.Errorf("swap error")
	}
}
