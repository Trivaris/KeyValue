package api

import (
	"crypto/rand"
	mrand "math/rand"
	"net"
	"net/netip"
	"slices"
	"sort"
	"sync"
	"time"
)

const (
    FlagMaster uint16 = 1 << iota
    FlagSlave  uint16 = 1 << iota
    FlagPfail  uint16 = 1 << iota
    FlagFail   uint16 = 1 << iota
    FlagMyself uint16 = 1 << iota
)

var crc16Table [256]uint16

type NodeID [20]byte

type Cluster struct {
	nodes []*Node
	myself *Node
	mu     sync.RWMutex
	node_timeout time.Duration
}

type Node struct {
	ID NodeID
	Addr netip.AddrPort
	CPort uint16
	Hostname string
	Flags uint16
	MasterID NodeID
	PingSent int64
	PongReceived int64
	ConfigEpoch uint64
	ReplOffset uint64
	LinkState bool
	Slots [2048]byte

	OutboundConn *net.TCPConn
	InboundConn *net.TCPConn
}

func NewCluster(hostname string, address netip.Addr, port uint16, cport uint16) (*Cluster, error) {
	var id [20]byte
	if _, err := rand.Read(id[:]); err != nil {
		return nil, err
	}

	myself := &Node{
		ID: id,
		Addr: netip.AddrPortFrom(address, port),
		CPort: cport,
		Hostname: hostname,
		Flags: FlagMyself | FlagMaster,
	}

	const poly = 0x1021 // CRC16-CCITT (XMODEM)
	for i := range 256 {
		crc := uint16(i) << 8
		for range 8 {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		crc16Table[i] = crc
	}

	return &Cluster{
		myself: myself,
	}, nil
}

func (c *Cluster) addNode(node *Node) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.nodes = append(c.nodes, node)
	return nil
}

// func (c *Cluster) getNode(id NodeID) *Node {
// 	i := sort.Search(len(c.nodes), func(i int) bool {
// 		return c.nodes[i].ID == id
// 	})
// 
// 	if i >= len(c.nodes) {
// 		return nil
// 	}
// 
// 	return c.nodes[i]
// }

func (c *Cluster) getNodeFromSlot(slot uint16) *Node {
	if slot >= 16384 {
		return nil
	}

	byteIndex := slot >> 3
	bitOffset := slot & 7
	mask := byte(1 << bitOffset)

	i := sort.Search(len(c.nodes), func(i int) bool {
		return (c.nodes[i].Slots[byteIndex] & mask) != 0
	})
	
	if i >= len(c.nodes) {
		return nil
	}

	return c.nodes[i]
}

func GetSlot(key string) uint16 {
	var crc uint16
	for i := range len(key) {
		crc = (crc << 8) ^ crc16Table[byte(crc>>8)^key[i]]
	}
	return crc & 0x3FFF
}

func (n *Node) isMyself() bool {
	return n.Flags & FlagMyself != 0
}

func (c *Cluster) getRandomStale() *Node {
	if len(c.nodes) == 0 {
		return nil
	}

	var picks [5]*Node
	for i, _ := range picks {
		randomIndex := mrand.Intn(len(c.nodes))
		picks[i] = c.nodes[randomIndex]
	}
	
	slices.SortFunc(picks[:], func(a *Node, b *Node) int {
		return int(a.PongReceived - b.PongReceived)
	})
	
	return picks[0]
}

func (c *Cluster) getStalest() *Node {
	if len(c.nodes) == 0 {
		return nil
	}

	sorted := c.nodes
	slices.SortFunc(sorted, func(a *Node, b *Node) int {
		return int(a.PongReceived - b.PongReceived)
	})

	return sorted[0]
}