package utils

import "context"

type Iterable[E any] interface {

	// Iterator use for range
	// it very slow, not recommend using it
	// example:
	//
	//  iter := iterable.Iterator()
	// 	for e := range iter.C {
	// 		do something...
	//	}
	Iterator() *Iterator[E]

	// Iter use for range
	// it very slow, not recommend using it
	// example:
	//
	// 	for e := range iterable.Iter() {
	// 		do something...
	//	}
	Iter() <-chan E

	// Foreach range the container
	// return false to break loop
	Foreach(func(index int, element E) bool)
}

// Iterator defines an iterator over a Set, its C channel can be used to range over the Set's
// elements.
type Iterator[E any] struct {
	C      <-chan E
	ctx    context.Context
	cancel context.CancelFunc
}

func (i *Iterator[E]) Ctx() context.Context {
	return i.ctx
}

// Stop stops the Iterator, no further elements will be received on C, C will be closed.
func (i *Iterator[E]) Stop() {
	i.cancel()

	// Exhaust any remaining elements.
	for range i.C {
	}
}

// newIterator returns a new Iterator instance together with its item and stop channels.
func newIterator[E any]() (*Iterator[E], chan<- E) {
	itemChan := make(chan E)
	ctx, cancel := context.WithCancel(context.Background())
	return &Iterator[E]{
		C:      itemChan,
		ctx:    ctx,
		cancel: cancel,
	}, itemChan
}
