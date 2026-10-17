package main

import(
	"os"
	"os/signal"
	"fmt"
	"syscall"
	"context"
	"git.trivaris.org/trivaris/keyvalue/cmd/internal/resp"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		dataMan 	:= DataManagerImpl{
			values: map[string]resp.RespItem{
				"foo":resp.RespString{Value: "bar"},
				"greeting":resp.RespString{Value: "Hello world!"},
			},
		}
		logicMan 	:= LogicManagerImpl{ dataMan: &dataMan }
		netManager 	:= NetManagerImpl{ logicMan: logicMan }
		netManager.Listen()
	}()

	<-ctx.Done()
	fmt.Println()
    fmt.Println(context.Cause(ctx))
    fmt.Println("exiting")

}
