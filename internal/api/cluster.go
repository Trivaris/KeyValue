package api

import (
	"fmt"
	"hash/fnv"
	"net"
	"sort"
	"sync"

	"git.trivaris.org/trivaris/keyvalue/internal/config"
)

type Cluster struct {
	config *config.Config
	nodes  []*clusterNode
	mu     sync.RWMutex
}

type shard struct {
	conn *net.TCPConn
	name string
}

type clusterNode struct {
	id    uint32
	shard *shard
	mu    sync.RWMutex
}

func NewCluster(config *config.Config) *Cluster {
	selfShard := &shard{name: config.NodeName}
	selfNode := &clusterNode{id: ConsistentHash(config.NodeName), shard: selfShard}
	return &Cluster{
		config: config,
		nodes:  []*clusterNode{selfNode},
	}
}

func (c *Cluster) getNode(key string) *clusterNode {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keyID := ConsistentHash(key)
	idx := sort.Search(len(c.nodes), func(i int) bool {
		return c.nodes[i].id >= keyID
	})
	if idx >= len(c.nodes) {
		return c.nodes[0]
	}
	return c.nodes[idx]
}

func (c *Cluster) getNodeID(key string) (uint32, bool) {
	node := c.getNode(key)

	node.mu.RLock()
	defer node.mu.RUnlock()

	return node.id, node.shard.name == c.config.NodeName
}

func (c *Cluster) getShardConn(nodeID uint32) (*net.TCPConn, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	idx := sort.Search(len(c.nodes), func(i int) bool {
		return c.nodes[i].id == nodeID
	})
	if idx >= len(c.nodes) {
		return nil, false
	}

	node := c.nodes[idx]

	node.mu.RLock()
	defer node.mu.RUnlock()

	return node.shard.conn, true
}

func (c *Cluster) addShardConn(name string, conn *net.TCPConn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Printf("adding %s with connection on %s\n", name, conn.RemoteAddr().String())

	id := ConsistentHash(name)
	shard := &shard{conn: conn, name: name}
	node := &clusterNode{id: id, shard: shard}
	c.nodes = append(c.nodes, node)
}

func ConsistentHash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	hash := h.Sum32()
	fmt.Printf("calculated id %d from key %q\n", hash, key)
	return hash
}
