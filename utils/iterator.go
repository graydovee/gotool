package utils

import "context"

type Iterable[E any] interface {
	Iterator() *Iterator[E]
	Iter() <-chan E
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

type Slice[E any] []E

func (s Slice[E]) Iter() <-chan E {
	ch := make(chan E)
	go func() {
		for _, elem := range s {
			ch <- elem
		}
		close(ch)
	}()
	return ch
}

func (s Slice[E]) Iterator() *Iterator[E] {
	iterator, ch := newIterator[E]()
	go func() {
		defer close(ch)
		for _, elem := range s {
			select {
			case <-iterator.Ctx().Done():
				return
			case ch <- elem:
			}
		}
	}()
	return iterator
}
