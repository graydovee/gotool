package utils

import "sync"

type Map[Key comparable, Value any] struct {
	syncMap    sync.Map
	emptyValue Value
}

func (m *Map[Key, Value]) Load(k Key) (Value, bool) {
	load, ok := m.syncMap.Load(k)
	if ok {
		return load.(Value), true
	}
	return m.emptyValue, false
}

func (m *Map[Key, Value]) Store(k Key, v Value) {
	m.syncMap.Store(k, v)
}

func (m *Map[Key, Value]) LoadOrStore(k Key, v Value) (Value, bool) {
	store, loaded := m.syncMap.LoadOrStore(k, v)
	if store == nil {
		return m.emptyValue, loaded
	}
	return store.(Value), loaded
}

func (m *Map[Key, Value]) LoadAndDelete(k Key) (Value, bool) {
	value, deleted := m.syncMap.LoadAndDelete(k)
	if deleted {
		return value.(Value), true
	}
	return m.emptyValue, false
}

func (m *Map[Key, Value]) Delete(k Key) {
	m.LoadAndDelete(k)
}

func (m *Map[Key, Value]) Range(f func(k Key, v Value) bool) {
	m.syncMap.Range(func(key, value any) bool {
		return f(key.(Key), value.(Value))
	})
}
