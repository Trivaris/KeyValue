package api

type ClusterBusHandler struct {
	store *Store
	cluster *Cluster
}

func NewClusterBusHandler(store *Store, cluster *Cluster) *ClusterBusHandler {
	return &ClusterBusHandler{
		store: store,
		cluster: cluster,
	}
}

func (cb *ClusterBusHandler) handleRequest(data []byte) []byte {
	switch data[0] {
	case '0': {
		// Ping packet
	}
	case '1':  {
		// Pong packet
	}
	case '2': {
		// Meet packet
	}
	}

	return nil
}