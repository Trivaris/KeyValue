package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.trivaris.org/trivaris/keyvalue/internal/api"
	"git.trivaris.org/trivaris/keyvalue/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if len(os.Args) < 2 {
			fmt.Println("No Node Name specified")
			return
		}

		config := config.NewConfig(os.Args[1:])
		cluster := api.NewCluster(config)
		store := api.NewStore(cluster)
		handler := api.NewHandler(store, cluster, config)
		server := api.NewServer(handler, config)

		server.Listen()
	}()

	<-ctx.Done()
	fmt.Println()
	fmt.Println(context.Cause(ctx))
	fmt.Println("exiting")

}
