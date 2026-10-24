package main

import (
	"flag"
	"fmt"
	"strconv"

	"git.trivaris.org/trivaris/keyvalue/internal/core"
	"git.trivaris.org/trivaris/keyvalue/internal/core/command"
	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

func main() {
	commandType := flag.String("type", "GET", "Command Type")
	commandAddress := flag.String("address", "127.0.0.1", "Target Address")
	commandPort := flag.Int("port", 6379, "Target Port")

	var commandArgs []serialization.Entity

	flag.Func("rawArg", "Raw Command Arg", func(val string) error {
		entity, _, err := serialization.Deserialize([]byte(val), 0)
		if err != nil {
			return err
		}
    	commandArgs = append(commandArgs, entity)
    	return nil
	})

	flag.Func("bulkArg", "Bulk String Arg", func(val string) error {
    	commandArgs = append(commandArgs, serialization.BulkString{ Value: val })
    	return nil
	})

	flag.Func("errorArg", "Error Arg", func(val string) error {
    	commandArgs = append(commandArgs, serialization.Error{ Error: fmt.Errorf(val) })
    	return nil
	})

	flag.Func("intArg", "Integer Arg", func(val string) error {
    	n, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		commandArgs = append(commandArgs, serialization.Integer{ Value: int64(n) })
    	return nil
	})

	flag.Func("nullArg", "Null Arg", func(_ string) error {
    	commandArgs = append(commandArgs, serialization.Null{})
    	return nil
	})

	flag.Parse()

	fmt.Printf("sending %s request to %s:%d, with args %v\n", *commandType, *commandAddress, *commandPort, commandArgs)

	cmd, err := command.Parse(*commandType, commandArgs)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := core.CreateConn(fmt.Sprintf("%s:%d", *commandAddress, *commandPort))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = conn.Close() }()

	response, err := command.Execute(cmd, conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(response)
}
