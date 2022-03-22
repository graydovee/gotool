package utils

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

func (s Slice[E]) Foreach(f func(index int, element E) bool) {
	for i, e := range s {
		if !f(i, e) {
			break
		}
	}
}

func (s Slice[E]) Swap(i, j int) {
	Swap(s, i, j)
}

func Swap[S ~[]E, E any](arr S, i, j int) {
	t := arr[i]
	arr[i] = arr[j]
	arr[j] = t
}
