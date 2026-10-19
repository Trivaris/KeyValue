package main

import (
	"fmt"
	"sync"

	"git.trivaris.org/trivaris/keyvalue/internal/resp"
)

const DATA_INFO = "[DATA INFO]"

type store struct {
	cluster *cluster
	values  map[string]resp.RespItem
	client  *client
	mu      sync.RWMutex
}

func (st *store) getValue(key string) (resp.RespItem, bool) {
	id, local := st.cluster.getNodeID(key)
	if local {
		fmt.Printf("%s Getting %q from Local Stoarge, keyId: %d\n", DATA_INFO, key, id)
		return st.getLocal(key)
	} else {
		fmt.Printf("%s Getting %q from Cloud Stoarge, keyId: %d\n", DATA_INFO, key, id)
		return st.getFromNode(key, id)
	}
}

func (st *store) setValue(key string, value resp.RespItem) error {
	id, local := st.cluster.getNodeID(key)
	if local {
		fmt.Printf("%s Setting %q to %#v in Local Stoarge, keyId: %d\n", DATA_INFO, key, value, id)
		st.setLocal(key, value)
		return nil
	} else {
		fmt.Printf("%s Setting %q to %#v in Cloud Stoarge, keyId: %d\n", DATA_INFO, key, value, id)
		return st.setOnNode(key, value, id)
	}
}

func (st *store) getLocal(key string) (resp.RespItem, bool) {
	st.mu.RLocker().Lock()
	defer st.mu.RLocker().Unlock()

	item, ok := st.values[key]
	return item, ok
}

func (st *store) getFromNode(key string, id uint32) (resp.RespItem, bool) {
	conn, ok := st.cluster.getShardConn(id)
	if !ok {
		return nil, false
	}
	request := resp.RespArray{Values: []resp.RespItem{
		resp.RespString{Value: "GET"},
		resp.RespString{Value: key},
	}}.Serialize()

	item, err := st.client.sendRequest(request, conn)
	if err != nil {
		return nil, false
	}

	arr, ok := item.(resp.RespArray)
	if !ok || len(arr.Values) == 0 {
		return nil, false
	}
	return arr.Values[0], true
}

func (st *store) setLocal(key string, value resp.RespItem) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.values[key] = value
}

func (st *store) setOnNode(key string, value resp.RespItem, id uint32) error {
	conn, ok := st.cluster.getShardConn(id)
	if !ok {
		return nil
	}

	request := resp.RespArray{Values: []resp.RespItem{
		resp.RespString{Value: "SET"},
		resp.RespString{Value: key},
		value,
	}}.Serialize()

	item, err := st.client.sendRequest(request, conn)
	if err != nil {
		return err
	}
	if !resp.IsSuccessString(item) {
		return fmt.Errorf("Not a successstring")
	}

	return nil
}
