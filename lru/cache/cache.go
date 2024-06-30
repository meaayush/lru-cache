package cache

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	key, value interface{}
	prev, next *Node
	expiration time.Time
}

type Cache struct {
	capacity   int
	cache      map[interface{}]*Node
	head, tail *Node
	lock       sync.RWMutex
}

func NewCache(capacity int) *Cache {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head
	return &Cache{
		capacity: capacity,
		cache:    make(map[interface{}]*Node),
		head:     head,
		tail:     tail,
	}
}

func (c *Cache) addNode(node *Node) {
	node.prev = c.head
	node.next = c.head.next
	c.head.next.prev = node
	c.head.next = node
}

func (c *Cache) removeNode(node *Node) {
	prev := node.prev
	next := node.next
	prev.next = next
	next.prev = prev
}

func (c *Cache) moveNodeToHead(node *Node) {
	c.removeNode(node)
	c.addNode(node)
}

func (c *Cache) popTail() *Node {
	res := c.tail.prev
	c.removeNode(res)
	return res
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if node, exists := c.cache[key]; exists {
		if time.Now().After(node.expiration) {
			c.removeNode(node)
			delete(c.cache, key)
			return nil, false
		}
		c.moveNodeToHead(node)
		return node.value, true
	}
	return nil, false
}

func (c *Cache) Set(key, value interface{}, ttl int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	expiration := time.Now().Add(time.Duration(ttl) * time.Second)

	if node, exists := c.cache[key]; exists {
		node.value = value
		node.expiration = expiration
		c.moveNodeToHead(node)
	} else {
		newNode := &Node{
			key:        key,
			value:      value,
			expiration: expiration,
		}
		c.cache[key] = newNode
		c.addNode(newNode)

		if len(c.cache) > c.capacity {
			tail := c.popTail()
			delete(c.cache, tail.key)
		}
	}
}

func (c *Cache) Delete(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if node, exists := c.cache[key]; exists {
		c.removeNode(node)
		delete(c.cache, key)
	}
}

func (c *Cache) GetAll() map[string]map[string]interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	result := make(map[string]map[string]interface{})
	now := time.Now()

	for key, node := range c.cache {
		if now.Before(node.expiration) {
			keyStr := fmt.Sprintf("%v", key)
			result[keyStr] = map[string]interface{}{
				"value":      node.value,
				"expiration": node.expiration,
			}
		}
	}

	return result
}

func (c *Cache) CleanUp() {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()

	for key, node := range c.cache {
		if now.After(node.expiration) {
			c.removeNode(node)
			delete(c.cache, key)
		}
	}
}

func (c *Cache) StartCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			c.CleanUp()
		}
	}()
}
