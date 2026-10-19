package main

import (
	"fmt"
	"hash/fnv"
	"net"
	"sort"
	"sync"
)

const SHARDING_INFO string = "[SHARDING INFO]"

type cluster struct {
	nodes []*clusterNode
	mu    sync.RWMutex
}

type shard struct {
	conn *net.TCPConn
	name string
	mu   sync.RWMutex
}

type clusterNode struct {
	id    uint32
	shard *shard
	mu    sync.RWMutex
}

func (c *cluster) getNode(key string) *clusterNode {
	keyID := ConsistentHash(key)
	idx := sort.Search(len(c.nodes), func(i int) bool {
		return c.nodes[i].id >= keyID
	})
	return c.nodes[idx]
}

func (c *cluster) getNodeID(key string) (uint32, bool) {
	node := c.getNode(key)
	return node.id, node.shard.name == NODE_NAME
}

func (c *cluster) getShardConn(nodeID uint32) (*net.TCPConn, bool) {
	idx := sort.Search(len(c.nodes), func(i int) bool {
		return c.nodes[i].id == nodeID
	})
	if idx >= len(c.nodes) {
		return nil, false
	}
	return c.nodes[idx].shard.conn, true
}

func (c *cluster) addShardConn(name string, conn *net.TCPConn) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := ConsistentHash(name)
	shard := &shard{conn: conn, name: name}
	node := &clusterNode{id: id, shard: shard}
	c.nodes = append(c.nodes, node)

	fmt.Printf("%s Adding %s with Connection on %s\n", SHARDING_INFO, name, conn.RemoteAddr().String())
	return false
}

func ConsistentHash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	hash := h.Sum32()
	fmt.Printf("%s Calculated id %d from key %q\n", SHARDING_INFO, hash, key)
	return hash
}
