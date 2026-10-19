package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"git.trivaris.org/trivaris/keyvalue/internal/core"
	"git.trivaris.org/trivaris/keyvalue/internal/core/command"
	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

var (
	commandType    string                   = "HELP"
	commandAddress string                   = "localhost"
	commandPort    uint16                   = 8080
	commandArgs    []serialization.RespItem = make([]serialization.RespItem, 0)
)

func main() {
	skip := 1
	for i := range len(os.Args) {
		if skip > 0 {
			skip--
			continue
		}

		switch os.Args[i] {
		case "-x":
			{
				if len(os.Args) <= i+1 {
					fmt.Println("No Type specified")
					return
				}
				commandType = os.Args[i+1]
				skip = 1
			}
		case "-a":
			{
				if len(os.Args) <= i+1 {
					fmt.Println("No Address specified")
					return
				}
				if strings.Contains(os.Args[i+1], ":") {
					fmt.Println("Declare port via -p")
					return
				}
				commandAddress = os.Args[i+1]
				skip = 1
			}
		case "-p":
			{
				if len(os.Args) <= i+1 {
					fmt.Println("No Port specified")
					return
				}
				port, err := strconv.Atoi(os.Args[i+1])
				if err != nil {
					fmt.Println(err)
					return
				}
				commandPort = uint16(port)
				skip = 1
			}
		case "-v":
			{
				for j, arg := range os.Args[i+1:] {
					if strings.HasPrefix(arg, "-") {
						skip = j
						break
					}

					if arg == "null" {
						commandArgs = append(commandArgs, serialization.RespNull{})
						continue
					}

					itemType, itemValue, found := strings.Cut(arg, ":")
					if !found {
						commandArgs = append(commandArgs, serialization.RespString{Value: arg})
						continue
					}

					switch strings.ToLower(itemType) {
					case "string":
						commandArgs = append(commandArgs, serialization.RespString{Value: itemValue})
					case "int":
						{
							num, err := strconv.Atoi(itemValue)
							if err != nil {
								fmt.Println(err)
								return
							}
							commandArgs = append(commandArgs, serialization.RespInteger{Value: int64(num)})
						}
					default:
						{
							commandArgs = append(commandArgs, serialization.RespString{Value: arg})
						}
					}
				}
			}
		}
	}

	fmt.Printf("sending %s request to %s:%d, with args %v\n", commandType, commandAddress, commandPort, commandArgs)

	cmd, err := command.Parse(commandType, commandArgs)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := core.CreateConn(fmt.Sprintf("%s:%d", commandAddress, commandPort))
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
