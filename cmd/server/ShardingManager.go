package main

import (
	"fmt"
	"hash/fnv"
	"net"
	"sync"
)

const SHARDING_INFO string = "[SHARDING INFO]"

type shardingManager struct {
	shards []shard
	headNode *llNode
	mu sync.RWMutex
}

type shard struct {
	conn *net.TCPConn
	mu sync.RWMutex
}

type llNode struct {
	next *llNode
	id uint32
	shard *shard
	mu sync.RWMutex
}

// Returns id, isHere
func (this *shardingManager) getServerId(key string) (uint32, bool) {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if this.headNode == nil {
		return 0, false
	}

	hash := ConsistentHash(key)
	cur := this.headNode
	for {
		if cur.id < cur.next.id {
			if cur.id < hash && hash <= cur.next.id {
				return cur.next.id, cur.next == this.headNode
			}
		} else {
			if hash > cur.id || hash <= cur.next.id {
				return cur.next.id, cur.next == this.headNode
			}
		}

		cur = cur.next
		if cur == this.headNode {
			return this.headNode.id, true
		}
	}
}

func ConsistentHash(key string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(key))
    hash := h.Sum32()
	fmt.Printf("%s Calculated id %d from key %q\n", SHARDING_INFO, hash, key)
	return hash
}

func (this *shardingManager) getServerConn(id uint32) (*net.TCPConn, bool) {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if this.headNode == nil {
		return nil, false
	}
	cur := this.headNode
	for {
		if cur.id == id {
			if cur.shard == nil {
				return nil, false
			}
			return cur.shard.conn, true
		}
		cur = cur.next
		if cur == this.headNode {
			return nil, false
		}
	}
}

func (this *shardingManager) addServerConn(name string, conn *net.TCPConn) bool {
	this.mu.Lock()
	defer this.mu.Unlock()

	id := ConsistentHash(name)
	shard := &shard{conn: conn}
	fmt.Printf("%s Adding %s with Connection on %s\n", SHARDING_INFO, name, conn.RemoteAddr().String())
	return this.headNode.insertSortedById(id, shard)
}

func (this *llNode) insertSortedById(id uint32, shard *shard) bool {
	this.mu.Lock()
	defer this.mu.Unlock()

	cur := this
	for {
		if cur.id == id {
			return true
		}

		isBetween := cur.id < cur.next.id && (cur.id < id && id < cur.next.id)
		isAtSeam := cur.id >= cur.next.id && (id > cur.id || id < cur.next.id)

		if isBetween || isAtSeam {
            newNode := &llNode{
                next:  cur.next,
                id:    id,
                shard: shard,
            }
            cur.next = newNode
            return false
        }
		cur = cur.next
	}
}
