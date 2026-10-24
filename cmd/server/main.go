package main

import (
	"context"
	"flag"
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
		listenAddress := flag.String("address", "0.0.0.0", "Listen Address")
		listenPort := flag.Int("port", 6379, "Listen Port")
		cListenPort := flag.Int("cport", 16379, "Cluster Bus Listen Port")
		hostname := flag.String("hostname", "node1", "Hostname")
		logPath := flag.String("logPath", "./log.txt", "Path to transaction Log File")


		config, err := config.NewNodeConfig(*hostname, *listenAddress, uint16(*listenPort), uint16(*cListenPort))
		if err != nil {
			fmt.Println(err)
			return
		}

		log, err := api.NewTransactionLog(*logPath)		
		if err != nil {
			fmt.Println(err)
			return
		}
		
		defer log.Close()

		cluster := api.NewCluster(config)
		store := api.NewStore(cluster)
		handler := api.NewHandler(store, cluster, log, config)
		server := api.NewServer(handler, config)

		server.Listen()
	}()

	<-ctx.Done()
	fmt.Println()
	fmt.Println(context.Cause(ctx))
	fmt.Println("exiting")

}
