package utils

import (
	"fmt"
	"testing"
)

type TypeA struct {
	v int
}

func TestMap(t *testing.T) {
	var m Map[string, TypeA]
	m.Store("1", TypeA{1})

	a, _ := m.Load("1")
	b, _ := m.Load("2")
	fmt.Println(a, b)
}

func TestMapI(t *testing.T) {
	var m Map[string, *TypeA]
	m.Store("1", &TypeA{1})

	a, _ := m.Load("1")
	b, _ := m.Load("2")
	fmt.Println(a, b)
}
