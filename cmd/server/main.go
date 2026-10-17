package main

import(
	"os"
	"os/signal"
	"fmt"
	"syscall"
	"context"
	"git.trivaris.org/trivaris/keyvalue/cmd/internal/resp"
)

var NODE_NAME string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if len(os.Args) < 2 {
			fmt.Println("No Node Name specified")
			return
		}
		NODE_NAME = os.Args[1]
		values := map[string]resp.RespItem{
			"foo":		resp.RespString{Value: "bar"},
			"greeting":	resp.RespString{Value: "Hello world!"},
		}

		headNode := &llNode			{ next: nil, id: ConsistentHash(NODE_NAME), shard: nil }
		headNode.next = headNode
		sharding := shardingManager { shards: make([]shard, 0), headNode: headNode }
		data 	 := dataManager     { values: values, sharding: &sharding }
		logic 	 := logicManager    { data:   &data,  sharding: &sharding }
		net 	 := netManager      { logic:  &logic }

		net.listen()
	}()

	<-ctx.Done()
	fmt.Println()
    fmt.Println(context.Cause(ctx))
    fmt.Println("exiting")

}
