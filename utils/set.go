package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type empty struct{}

type Set[E comparable] map[E]empty

// An OrderedPair represents a 2-tuple of values.
type OrderedPair[E comparable] struct {
	First  E
	Second E
}

func NewSet[E comparable]() Set[E] {
	return make(Set[E])
}

// Equal says whether two 2-tuples contain the same values in the same order.
func (pair *OrderedPair[E]) Equal(other OrderedPair[E]) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}

	return false
}

func (set *Set[E]) Add(i E) bool {
	_, found := (*set)[i]
	if found {
		return false //False if it existed already
	}

	(*set)[i] = empty{}
	return true
}

func (set *Set[E]) Contains(i ...E) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}
	return true
}

func (set *Set[E]) IsSubset(other Set[E]) bool {
	if set.Len() > other.Len() {
		return false
	}
	for elem := range *set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set *Set[E]) IsProperSubset(other Set[E]) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *Set[E]) IsSuperset(other Set[E]) bool {
	return other.IsSubset(*set)
}

func (set *Set[E]) IsProperSuperset(other Set[E]) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *Set[E]) Union(other Set[E]) Set[E] {
	unionedSet := NewSet[E]()

	for elem := range *set {
		unionedSet.Add(elem)
	}
	for elem := range other {
		unionedSet.Add(elem)
	}
	return unionedSet
}

func (set *Set[E]) Intersect(other Set[E]) Set[E] {
	intersection := NewSet[E]()
	// loop over smaller set
	if set.Len() < other.Len() {
		for elem := range *set {
			if other.Contains(elem) {
				intersection.Add(elem)
			}
		}
	} else {
		for elem := range other {
			if set.Contains(elem) {
				intersection.Add(elem)
			}
		}
	}
	return intersection
}

func (set *Set[E]) Difference(other Set[E]) Set[E] {
	difference := NewSet[E]()
	for elem := range *set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return difference
}

func (set *Set[E]) SymmetricDifference(other Set[E]) Set[E] {
	aDiff := set.Difference(other)
	bDiff := other.Difference(*set)
	return aDiff.Union(bDiff)
}

func (set *Set[E]) Clear() {
	*set = NewSet[E]()
}

func (set *Set[E]) Remove(i E) {
	delete(*set, i)
}

func (set *Set[E]) Len() int {
	return len(*set)
}

func (set *Set[E]) Each(cb func(E) bool) {
	for elem := range *set {
		if !cb(elem) {
			break
		}
	}
}

func (set *Set[E]) Iter() <-chan E {
	ch := make(chan E)
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *Set[E]) Iterator() *Iterator[E] {
	iterator, ch := newIterator[E]()

	go func() {
		defer close(ch)
		for elem := range *set {
			select {
			case <-iterator.Ctx().Done():
				return
			case ch <- elem:
			}
		}

	}()

	return iterator
}

func (set *Set[E]) Equal(other Set[E]) bool {
	if set.Len() != other.Len() {
		return false
	}
	for elem := range *set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set *Set[E]) Clone() Set[E] {
	clonedSet := NewSet[E]()
	for elem := range *set {
		clonedSet.Add(elem)
	}
	return clonedSet
}

func (set *Set[E]) String() string {
	items := make([]string, 0, len(*set))

	for elem := range *set {
		items = append(items, fmt.Sprintf("%v", elem))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

// String outputs a 2-tuple in the form "(A, B)".
func (pair OrderedPair[E]) String() string {
	return fmt.Sprintf("(%v, %v)", pair.First, pair.Second)
}

func (set *Set[E]) Pop() interface{} {
	for item := range *set {
		delete(*set, item)
		return item
	}
	return nil
}

func (set *Set[E]) ToSlice() []E {
	keys := make([]E, 0, set.Len())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (set *Set[E]) MarshalJSON() ([]byte, error) {
	items := make([]string, 0, set.Len())

	for elem := range *set {
		b, err := json.Marshal(elem)
		if err != nil {
			return nil, err
		}

		items = append(items, string(b))
	}

	return []byte(fmt.Sprintf("[%s]", strings.Join(items, ","))), nil
}

// UnmarshalJSON recreates a set from a JSON array, it only decodes
// primitive types. Numbers are decoded as json.Number.
func (set *Set[E]) UnmarshalJSON(b []byte) error {
	var i []E

	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err := d.Decode(&i)
	if err != nil {
		return err
	}

	for _, v := range i {
		set.Add(v)
	}

	return nil
}
