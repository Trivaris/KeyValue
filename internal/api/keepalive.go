package api

import (
	"net"
	"net/netip"
	"time"
)

func KeepAlive(cluster *Cluster) {
	go func() {
		for {
			go pingRandom(cluster)
			go noticeNewConnections(cluster)
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		dur := time.Duration(float64(cluster.node_timeout) * 0.5)
		for {
			go pingStalest(cluster)
			time.Sleep(dur)
		}
	}()
}

func pingRandom(cluster *Cluster) {
	node := cluster.getRandomStale()
	sendPing(node)
}

func pingStalest(cluster *Cluster) {
	node := cluster.getStalest()
	sendPing(node)
}

func noticeNewConnections(cluster *Cluster) {
	for _, node := range cluster.nodes {
		if node.OutboundConn == nil {
			addr := netip.AddrPortFrom(node.Addr.Addr(), node.CPort)
			conn, err := net.DialTCP("tcp", nil, net.TCPAddrFromAddrPort(addr))
			if err != nil {
				// TODO: Implement Broadcast
				panic("node is down")
			}

			node.OutboundConn = conn
			go sendPing(node)
		}
	}
}

func sendPing(node *Node) {
	// TODO: Implement ping send
}