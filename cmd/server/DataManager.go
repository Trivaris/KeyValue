package main

import (
	"fmt"
	"sync"

	"git.trivaris.org/trivaris/keyvalue/cmd/internal/resp"
)

const DATA_INFO = "[DATA INFO]"

type dataManager struct {
	sharding *shardingManager
	values map[string]resp.RespItem
	mu sync.RWMutex
}

func (this *dataManager) getValue(key string) (resp.RespItem, bool) {
	id, isHere := this.sharding.getServerId(key)
	if isHere {
		fmt.Printf("%s Getting %q from Local Stoarge, keyId: %d\n", DATA_INFO, key, id)		
		return this.getLocal(key)
	}

	fmt.Printf("%s Getting %q from Cloud Stoarge, keyId: %d\n", DATA_INFO, key, id)		
	return this.getCloud(key, id)
}

func (this *dataManager) getLocal(key string) (resp.RespItem, bool) {
	this.mu.RLocker().Lock()
	defer this.mu.RLocker().Unlock()

	item, ok := this.values[key]
	return item, ok
}

func (this *dataManager) getCloud(key string, id uint32) (resp.RespItem, bool) {
	conn, ok := this.sharding.getServerConn(id)
	if !ok {
		return nil, false
	}
	request := resp.RespArray{Values: []resp.RespItem{
		resp.RespString{Value: "GET"},
		resp.RespString{Value: key},
	}}.Serialize()

	_, err := conn.Write(request)
	if err != nil {
		fmt.Printf("%s Write Error: %s\n", DATA_INFO, err)
		return nil, false
	}

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Read Error")
		return nil, false
	}

	item, _, err := resp.Deserialize(buf[:n], 0)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	arr, ok := item.(resp.RespArray)
	if !ok || len(arr.Values) == 0 {
		return nil, false
	}
	return arr.Values[0], true
}

func (this *dataManager) setValue(key string, value resp.RespItem) {
	id, isHere := this.sharding.getServerId(key)
	if isHere {
		fmt.Printf("%s Setting %q to %#v in Local Stoarge, keyId: %d\n", DATA_INFO, key, value, id)		
		this.setLocal(key, value)
		return
	}

	fmt.Printf("%s Setting %q to %#v in Cloud Stoarge, keyId: %d\n", DATA_INFO, key, value, id)		
	this.setCloud(key, value, id)
}

func (this *dataManager) setLocal(key string, value resp.RespItem) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.values[key] = value
}

func (this *dataManager) setCloud(key string, value resp.RespItem, id uint32) {
	conn, ok := this.sharding.getServerConn(id)
	if !ok {
		return
	}
	request := resp.RespArray{Values: []resp.RespItem{
		resp.RespString{Value: "SET"},
		resp.RespString{Value: key},
		value,
	}}.Serialize()

	_, err := conn.Write(request)
	if err != nil {
		fmt.Printf("Write Error")
		return
	}

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Read Error")
		return
	}

	item, _, err := resp.Deserialize(buf[:n], 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	response, ok := item.(resp.RespString)
	if !ok {
		fmt.Println(fmt.Errorf("CLOUD SET :: DID NOT RECEIVE STRING"))
		return
	}
	if response.Value != "Success!" {
		fmt.Println(fmt.Errorf("CLOUD SET :: DID NOT RECEIVE SUCCESS"))
		return
	}
}
