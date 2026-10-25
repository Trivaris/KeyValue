package api

import (
	"fmt"
	"sync"

	"git.trivaris.org/trivaris/keyvalue/internal/core"
	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type Store struct {
	cluster *Cluster
	values  map[string]serialization.Entity
	mu      sync.RWMutex
}

func NewStore(cluster *Cluster) *Store {
	return &Store{
		cluster: cluster,
		values:  make(map[string]serialization.Entity),
	}
}

func (st *Store) getValue(key string) (serialization.Entity, bool) {
	slot := GetSlot(key)
	node := st.cluster.getNodeFromSlot(slot)
	if node.isMyself() {
		fmt.Printf("getting %q from local storage\n", key)
		return st.getLocal(key)
	} else {
		fmt.Printf("getting %q from node %d\n", key, node.ID[:])
		return st.getFromNode(key, node)
	}
}

func (st *Store) setValue(key string, value serialization.Entity) error {
	slot := GetSlot(key)
	node := st.cluster.getNodeFromSlot(slot)
	if node.isMyself() {
		fmt.Printf("setting %q to %#v in local storage\n", key, value)
		st.setLocal(key, value)
		return nil
	} else {
		fmt.Printf("setting %q to %#v on node %d\n", key, value, node.ID[:])
		return st.setOnNode(key, value, node)
	}
}

func (st *Store) getLocal(key string) (serialization.Entity, bool) {
	st.mu.RLocker().Lock()
	defer st.mu.RLocker().Unlock()

	item, ok := st.values[key]
	return item, ok
}

func (st *Store) getFromNode(key string, node *Node) (serialization.Entity, bool) {
	request := serialization.Array{Values: []serialization.Entity{
		serialization.BulkString{Value: "GET"},
		serialization.BulkString{Value: key},
	}}.Serialize()

	item, err := core.SendRequest(request, node.OutboundConn)
	if err != nil {
		return nil, false
	}

	arr, ok := item.(serialization.Array)
	if !ok || len(arr.Values) == 0 {
		return nil, false
	}
	return arr.Values[0], true
}

func (st *Store) setLocal(key string, value serialization.Entity) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.values[key] = value
}

func (st *Store) setOnNode(key string, value serialization.Entity, node *Node) error {
	request := serialization.Array{Values: []serialization.Entity{
		serialization.BulkString{Value: "SET"},
		serialization.BulkString{Value: key},
		value,
	}}.Serialize()

	item, err := core.SendRequest(request, node.OutboundConn)
	if err != nil {
		return err
	}

	if !serialization.IsSuccessString(item) {
		return fmt.Errorf("request unsuccessful")
	}

	return nil
}
