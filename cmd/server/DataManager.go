package main

import (
	"sync"
	"git.trivaris.org/trivaris/keyvalue/cmd/internal/resp"
)

type dataManager struct {
	values map[string]resp.RespItem
	mu sync.RWMutex
}

func (this *dataManager) getValue(key string) (resp.RespItem, bool) {
	this.mu.RLocker().Lock()
	defer this.mu.RLocker().Unlock()

	item, ok := this.values[key]
	return item, ok
}

func (this *dataManager) setValue(key string, value resp.RespItem) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.values[key] = value
}
