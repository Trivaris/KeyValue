package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
			"foo":      resp.RespString{Value: "bar"},
			"greeting": resp.RespString{Value: "Hello world!"},
		}

		selfShard 	:= &shard		{ name: NODE_NAME }
		selfNode 	:= &clusterNode	{ id: ConsistentHash(NODE_NAME), shard: selfShard }

		cluster 	:= &cluster		{ nodes: []*clusterNode{ selfNode } }
		data 		:= &store		{ values: values, cluster: cluster }
		handler 	:= &handler		{ data: data, cluster: cluster }
		server 		:= &server		{ handler: handler }

		server.listen()
	}()

	<-ctx.Done()
	fmt.Println()
	fmt.Println(context.Cause(ctx))
	fmt.Println("exiting")

}
