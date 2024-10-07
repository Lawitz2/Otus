package hw04lrucache

type Key string

/*
if ItemKey from ListItem (list.go) is removed:
- comment lines 27, 32, 42, 50
- uncomment the loop at lines 33-38
*/

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, val interface{}) bool {
	if _, ok := l.items[key]; ok {
		l.items[key].Value = val
		l.queue.MoveToFront(l.items[key])
		l.queue.Front().ItemKey = key
		l.items[key] = l.queue.Front()
		return true
	}
	if l.queue.Len() == l.capacity {
		delete(l.items, l.queue.Back().ItemKey)
		// for mapKey, item := range l.items {
		//	if l.queue.Back() == item {
		//		delete(l.items, mapKey)
		//		break
		//	}
		// }
		l.queue.Remove(l.queue.Back())
	}
	l.queue.PushFront(val)
	l.queue.Front().ItemKey = key
	l.items[key] = l.queue.Front()
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(l.items[key])
		l.queue.Front().ItemKey = key
		l.items[key] = l.queue.Front()
		return item.Value, ok
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
