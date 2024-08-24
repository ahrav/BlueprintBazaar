package cache

import (
	"errors"
	"fmt"
)

type DoubleLinkNode struct {
	Key  int
	Val  any
	Next *DoubleLinkNode
	Prev *DoubleLinkNode
}

type DoublyLinkedList struct {
	head *DoubleLinkNode
	tail *DoubleLinkNode
	size int
}

func (dll *DoublyLinkedList) RemoveNode(node *DoubleLinkNode) {
	if node == nil {
		return
	}

	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		dll.head = node.Next
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		dll.tail = node.Prev
	}

	node.Prev = nil
	node.Next = nil

	dll.size--
}

func (dll *DoublyLinkedList) AddItemToHead(item any) *DoubleLinkNode {
	node := &DoubleLinkNode{Val: item}
	dll.AddToHead(node)
	return node
}

func (dll *DoublyLinkedList) AddToHead(node *DoubleLinkNode) {
	if dll.head == nil {
		dll.head = node
		dll.tail = node
	} else {
		node.Next = dll.head
		dll.head.Prev = node
		dll.head = node
	}
}

func (dll *DoublyLinkedList) RemoveTail() *DoubleLinkNode {
	if dll.tail == nil {
		return nil
	}

	prevTail := dll.tail
	if dll.tail.Prev != nil {
		dll.tail = dll.tail.Prev
		dll.tail.Next = nil
	} else {
		dll.head = nil
		dll.tail = nil
	}
	prevTail.Prev = nil

	dll.size--

	return prevTail
}

type MemoryStorage struct {
	store  map[any]any
	cap    int
	size   int
	isFull bool
}

func NewMemoryStorage(cap int) *MemoryStorage {
	return &MemoryStorage{
		cap:   cap,
		store: make(map[any]any, cap),
	}
}

var ErrStorageFull = errors.New("storage is full")

func (m *MemoryStorage) Add(key any, value any) error {
	if m.isFull {
		return ErrStorageFull
	}

	m.size++
	if m.size == m.cap {
		m.isFull = true
	}
	m.store[key] = value

	return nil
}

func (m *MemoryStorage) Remove(key any) {
	delete(m.store, key)
}

func (m *MemoryStorage) Get(key any) (any, bool) {
	val, ok := m.store[key]
	return val, ok
}

type Storage interface {
	Add(key any, value any) error
	Remove(key any)
	Get(key any) (any, bool)
}

type LRUEvictionPolicy struct {
	dll *DoublyLinkedList
	lut map[any]*DoubleLinkNode
}

func NewLRUEvictionPolicy() *LRUEvictionPolicy {
	return &LRUEvictionPolicy{
		dll: new(DoublyLinkedList),
		lut: make(map[any]*DoubleLinkNode),
	}
}

func (l *LRUEvictionPolicy) ItemAccessed(item any) {
	if node, ok := l.lut[item]; ok {
		l.dll.RemoveNode(node)
		l.dll.AddToHead(node)
	} else {
		node := l.dll.AddItemToHead(item)
		l.lut[item] = node
	}
}

func (l *LRUEvictionPolicy) Evict(item any) any {
	return l.dll.RemoveTail()
}

type EvictionPolicy interface {
	ItemAccessed(item any)
	Evict() any
}

type Cache struct {
	store  Storage
	policy EvictionPolicy
}

func NewCache(storage Storage, policy EvictionPolicy) *Cache {
	cache := Cache{
		store:  storage,
		policy: policy,
	}

	return &cache
}

func (c *Cache) Get(key any) (any, error) {
	val, ok := c.store.Get(key)
	if !ok {
		return nil, fmt.Errorf("key doesn't exist")
	}
	c.policy.ItemAccessed(key)

	return val, nil
}

func (c *Cache) Put(key, value any) error {
	c.policy.ItemAccessed(key)
	err := c.store.Add(key, value)
	if err != nil {
		if errors.Is(err, ErrStorageFull) {
			c.policy.Evict()
			if err = c.store.Add(key, value); err != nil {
				return fmt.Errorf("unkown error")
			}
		}
	}

	return nil
}
