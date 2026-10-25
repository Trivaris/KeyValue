package api

import (
	"fmt"

	"git.trivaris.org/trivaris/keyvalue/internal/core"
	"git.trivaris.org/trivaris/keyvalue/internal/core/command"
	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type Handler struct {
	store   *Store
	cluster *Cluster
	log     *TransactionLog
}

func NewHandler(store *Store, cluster *Cluster, log *TransactionLog) *Handler {
	return &Handler{
		store:   store,
		cluster: cluster,
		log:     log,
	}
}

func (h *Handler) handleRequest(buf []byte) serialization.Entity {
	item, _, err := serialization.Deserialize(buf, 0)

	if err != nil {
		fmt.Println(err)
		return serialization.Error{Error: err}
	}

	arr, ok := item.(serialization.Array)
	if !ok {
		err := fmt.Errorf("no command array found in %s", item.Serialize())
		fmt.Println(err)
		return serialization.Error{Error: err}
	}

	cmd, err := command.ParseRaw(arr)
	if err != nil {
		fmt.Println(err)
		return serialization.Error{Error: err}
	}

	switch c := cmd.(type) {
	case command.Get:
		return h.handleGet(c)
	case command.Set:
		return h.handleSet(c, arr)
	case command.ClusterMeet:
		return h.handleClusterMeet(c)
	default:
		{
			err := fmt.Errorf("request not yet implemented: %v", arr.Values)
			fmt.Println(err)
			return serialization.Error{Error: err}
		}
	}
}

func (h *Handler) handleGet(cmd command.Get) serialization.Entity {
	values := make([]serialization.Entity, len(cmd.Keys))

	for i, key := range cmd.Keys {
		value, ok := h.store.getValue(key)
		if !ok {
			value = serialization.Null{}
		}
		values[i] = value
	}

	return serialization.Array{Values: values}
}

func (h *Handler) handleSet(cmd command.Set, raw serialization.Array) serialization.Entity {
	h.log.Log(raw)
	err := h.store.setValue(cmd.Key, cmd.Value)

	if err != nil {
		fmt.Println(err)
		return serialization.Error{Error: err}
	}

	return serialization.SuccessString()
}

func (h *Handler) handleClusterMeet(cmd command.ClusterMeet) serialization.Entity {
	dialConn, err := core.CreateConn(cmd.Address)

	if err != nil {
		fmt.Println(err)
		return serialization.Error{Error: err}
	}

	request := command.Ping{Name: string(h.cluster.myself.ID[:]), Port: uint16(h.cluster.myself.Addr.Port())}

	item, err := command.Execute(request, dialConn)
	if err != nil {
		fmt.Println(err)
		return serialization.Error{Error: err}
	}

	if !serialization.IsSuccessString(item) {
		return serialization.Error{Error: fmt.Errorf("PING Request Unsuccesful")}
	}

	// TODO: Fixup, this is not correctly implemented
	h.cluster.addNode(&Node{
		OutboundConn: dialConn,
	})

	return serialization.SuccessString()
}
