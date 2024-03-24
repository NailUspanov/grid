package main

import "sync"

type ConcurrentMap struct {
	m  map[string]Node
	mu sync.Mutex
}

func NewConcurrentMap(cm map[string]Node) *ConcurrentMap {
	return &ConcurrentMap{
		m: cm,
	}
}

func (cm *ConcurrentMap) Set(key string, value Node) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.m[key] = value
}

func (cm *ConcurrentMap) Get(key string) (Node, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	value, ok := cm.m[key]
	return value, ok
}

func (cm *ConcurrentMap) Delete(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.m, key)
}
