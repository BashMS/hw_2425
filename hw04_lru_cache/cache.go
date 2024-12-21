package hw04lrucache

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
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	tempCacheItem := cacheItem{
		key:   key,
		value: value,
	}
	var isExist bool
	if item, ok := l.items[key]; ok {
		item.Value = tempCacheItem
		l.queue.MoveToFront(item)
		isExist = true
		return isExist
	} else {
		item := l.queue.PushFront(tempCacheItem)
		l.items[key] = item
		if l.queue.Len() > l.capacity {
			lastItem := l.queue.Back()
			l.queue.Remove(lastItem)
			delete(l.items, lastItem.Value.(cacheItem).key)
		}
		isExist = false
		return isExist
	}
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	var isExist bool
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		isExist = true
		return l.queue.Front().Value.(cacheItem).value, isExist
	}
	isExist = false
	return nil, isExist
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
