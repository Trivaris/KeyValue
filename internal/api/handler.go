package api

import (
	"fmt"
	"net"

	"git.trivaris.org/trivaris/keyvalue/internal/config"
	"git.trivaris.org/trivaris/keyvalue/internal/core"
	"git.trivaris.org/trivaris/keyvalue/internal/core/command"
	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type Handler struct {
	store   *Store
	cluster *Cluster
	config  *config.Config
}

func NewHandler(store *Store, cluster *Cluster, config *config.Config) *Handler {
	return &Handler{
		store:   store,
		cluster: cluster,
		config:  config,
	}
}

func (h *Handler) handleRequest(buf []byte, conn *net.TCPConn) serialization.RespItem {
	item, _, err := serialization.Deserialize(buf, 0)

	if err != nil {
		fmt.Println(err)
		return serialization.RespError{Error: err}
	}

	arr, ok := item.(serialization.RespArray)
	if !ok {
		err := fmt.Errorf("no command array found in %s", item.Serialize())
		fmt.Println(err)
		return serialization.RespError{Error: err}
	}

	cmd, err := command.ParseRaw(arr)
	if err != nil {
		fmt.Println(err)
		return serialization.RespError{Error: err}
	}

	switch c := cmd.(type) {
	case command.RespGet:
		return h.handleGet(c)
	case command.RespSet:
		return h.handleSet(c)
	case command.RespAddServer:
		return h.handleAddServer(c)
	case command.RespPing:
		return h.handlePing(c, conn)
	default:
		{
			err := fmt.Errorf("request not yet implemented: %v", arr.Values)
			fmt.Println(err)
			return serialization.RespError{Error: err}
		}
	}
}

func (h *Handler) handleGet(cmd command.RespGet) serialization.RespItem {
	values := make([]serialization.RespItem, len(cmd.Keys))

	for i, key := range cmd.Keys {
		value, ok := h.store.getValue(key)
		if !ok {
			value = serialization.RespNull{}
		}
		values[i] = value
	}

	return serialization.RespArray{Values: values}
}

func (h *Handler) handleSet(cmd command.RespSet) serialization.RespItem {
	err := h.store.setValue(cmd.Key, cmd.Value)

	if err != nil {
		fmt.Println(err)
		return serialization.RespError{Error: err}
	}

	return serialization.SuccessString()
}

func (h *Handler) handleAddServer(cmd command.RespAddServer) serialization.RespItem {
	dialConn, err := core.CreateConn(cmd.Address)

	if err != nil {
		fmt.Println(err)
		return serialization.RespError{Error: err}
	}

	request := command.RespPing{Name: h.config.NodeName, Port: uint16(h.config.ListenPort)}

	item, err := command.Execute(request, dialConn)
	if err != nil {
		fmt.Println(err)
		return serialization.RespError{Error: err}
	}

	if !serialization.IsSuccessString(item) {
		return serialization.RespError{Error: fmt.Errorf("PING Request Unsuccesful")}
	}

	h.cluster.addShardConn(cmd.Name, dialConn)

	return serialization.SuccessString()
}

func (h *Handler) handlePing(cmd command.RespPing, conn *net.TCPConn) serialization.RespItem {
	tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		err := fmt.Errorf("could not convert to tcp addr")
		fmt.Println(err)
		return serialization.RespError{Error: err}
	}

	dialConn, err := core.CreateConn(fmt.Sprintf("%s:%d", tcpAddr.IP, cmd.Port))
	if err != nil {
		fmt.Println(err)
		return serialization.RespError{Error: err}
	}

	h.cluster.addShardConn(cmd.Name, dialConn)

	return serialization.SuccessString()
}
