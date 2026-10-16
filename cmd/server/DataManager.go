package main

import (
	"fmt"
	"sync"
	"git.trivaris.org/trivaris/keyvalue/cmd/internal/resp"
)

// Gets the Data
type DataManager interface {
	getValue(string) resp.RespItem
	setValue(string, resp.RespItem)
}

type DataManagerImpl struct {
	values map[string]resp.RespItem
	mu sync.RWMutex
}

func (dataMan *DataManagerImpl) getValue(key string) resp.RespItem {
	dataMan.mu.RLocker().Lock()
	defer dataMan.mu.RLocker().Unlock()

	item, ok := dataMan.values[key]
	if ok {
		return item
	}

	item, ok = getFromBackup(key)
	if !ok {
		return resp.RespNull{}
	}
	
	return item
}

func (dataMan *DataManagerImpl) setValue(key string, value resp.RespItem) {
	dataMan.mu.Lock()
	defer dataMan.mu.Unlock()

	dataMan.values[key] = value
}

func getFromBackup(key string) (resp.RespItem, bool) {
	return resp.RespString{Value: fmt.Sprintf("Got %q from DB!", key)}, true
}
