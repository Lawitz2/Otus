package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	lock     sync.RWMutex
}

func (l *lruCache) Set(key Key, val interface{}) bool { // adds new value (or updates existing one) to cache based on key
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, ok := l.items[key]; ok { // update value if the key already exists in cache
		l.items[key].value = val
		l.queue.MoveToFront(l.items[key])
		l.queue.Front().itemKey = key // reassign the itemkey to the moved node since it's a node with a different pointer
		l.items[key] = l.queue.Front()
		return true
	}
	if l.queue.Len() == l.capacity { // if cache is full - remove the least requested element
		delete(l.items, l.queue.Back().itemKey)
		l.queue.Remove(l.queue.Back())
	}
	l.queue.PushFront(val)        // add new element to cache
	l.queue.Front().itemKey = key // reassign the itemkey to the moved node since it's a node with a different pointer
	l.items[key] = l.queue.Front()
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) { // get a value from cache based on key
	l.lock.RLock()
	defer l.lock.RUnlock()
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(l.items[key])
		l.queue.Front().itemKey = key // reassign the itemkey to the moved node since it's a node with a different pointer
		l.items[key] = l.queue.Front()
		return item.value, ok
	}
	return nil, false
}

func (l *lruCache) Clear() { // fully clears the cache
	l.lock.Lock()
	defer l.lock.Unlock()
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		lock:     sync.RWMutex{},
	}
}
