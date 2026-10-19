package api

import (
	"fmt"
	"sync"

	"git.trivaris.org/trivaris/keyvalue/internal/core"
	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type Store struct {
	cluster *Cluster
	values  map[string]serialization.RespItem
	mu      sync.RWMutex
}

func NewStore(cluster *Cluster) *Store {
	return &Store{
		cluster: cluster,
		values:  make(map[string]serialization.RespItem),
	}
}

func (st *Store) getValue(key string) (serialization.RespItem, bool) {
	id, local := st.cluster.getNodeID(key)
	if local {
		fmt.Printf("getting %q from local storage\n", key)
		return st.getLocal(key)
	} else {
		fmt.Printf("getting %q from node %d\n", key, id)
		return st.getFromNode(key, id)
	}
}

func (st *Store) setValue(key string, value serialization.RespItem) error {
	id, local := st.cluster.getNodeID(key)
	if local {
		fmt.Printf("setting %q to %#v in local storage\n", key, value)
		st.setLocal(key, value)
		return nil
	} else {
		fmt.Printf("setting %q to %#v on node %d\n", key, value, id)
		return st.setOnNode(key, value, id)
	}
}

func (st *Store) getLocal(key string) (serialization.RespItem, bool) {
	st.mu.RLocker().Lock()
	defer st.mu.RLocker().Unlock()

	item, ok := st.values[key]
	return item, ok
}

func (st *Store) getFromNode(key string, id uint32) (serialization.RespItem, bool) {
	conn, ok := st.cluster.getShardConn(id)
	if !ok {
		return nil, false
	}
	request := serialization.RespArray{Values: []serialization.RespItem{
		serialization.RespString{Value: "GET"},
		serialization.RespString{Value: key},
	}}.Serialize()

	item, err := core.SendRequest(request, conn)
	if err != nil {
		return nil, false
	}

	arr, ok := item.(serialization.RespArray)
	if !ok || len(arr.Values) == 0 {
		return nil, false
	}
	return arr.Values[0], true
}

func (st *Store) setLocal(key string, value serialization.RespItem) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.values[key] = value
}

func (st *Store) setOnNode(key string, value serialization.RespItem, id uint32) error {
	conn, ok := st.cluster.getShardConn(id)
	if !ok {
		return nil
	}

	request := serialization.RespArray{Values: []serialization.RespItem{
		serialization.RespString{Value: "SET"},
		serialization.RespString{Value: key},
		value,
	}}.Serialize()

	item, err := core.SendRequest(request, conn)
	if err != nil {
		return err
	}
	if !serialization.IsSuccessString(item) {
		return fmt.Errorf("connection unsuccessful")
	}

	return nil
}
