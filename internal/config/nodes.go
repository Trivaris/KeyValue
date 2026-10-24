package config

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"net"
	"net/netip"
	"slices"
	"strconv"
	"strings"
)

const (
    FlagMaster uint16 = 1 << iota
    FlagSlave  uint16 = 1 << iota
    FlagPfail  uint16 = 1 << iota
    FlagFail   uint16 = 1 << iota
    FlagMyself uint16 = 1 << iota
)

type NodeID [20]byte

type NodeConfig struct {
	Nodes []*Node
	Myself *Node
	counter uint64
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

func NewNodeConfig(hostname string, address string, port uint16, cport uint16) (*NodeConfig, error) {
	var id [20]byte
	if _, err := rand.Read(id[:]); err != nil {
		return nil, err
	}

	var addr [4]byte
	parts := strings.Split(":", address)
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid ip address")
	}

	for i := range addr {
		n, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, err
		}
		addr[i] = uint8(n)
	}

	myself := &Node{
		ID: id,
		Addr: netip.AddrPortFrom(netip.AddrFrom4(addr), port),
		CPort: cport,
		Hostname: hostname,
		Flags: FlagMyself | FlagMaster,
	}

	return &NodeConfig{
		Myself: myself,
		counter: 1,
	}, nil
}

func (c *NodeConfig) getRandomStale() *Node {
	if len(c.Nodes) == 0 {
		return nil
	}

	var picks [5]*Node
	for i, _ := range picks {
		randomIndex := mrand.Intn(len(c.Nodes))
		picks[i] = c.Nodes[randomIndex]
	}
	
	slices.SortFunc(picks[:], func(a *Node, b *Node) int {
		return int(a.PongReceived - b.PongReceived)
	})
	
	return picks[0]
}

func (c *NodeConfig) getStalest() *Node {
	if len(c.Nodes) == 0 {
		return nil
	}

	sorted := c.Nodes
	slices.SortFunc(sorted, func(a *Node, b *Node) int {
		return int(a.PongReceived - b.PongReceived)
	})

	return sorted[0]
}