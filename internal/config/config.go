package config

import "strconv"

type Config struct {
	NodeName      string
	ListenPort    uint16
	ListenAddress string
}

func NewConfig(args []string) *Config {
	defaultArgs := []string{"node1", "0.0.0.0", "8080"}
	for i, item := range args {
		if i >= len(defaultArgs) {
			break
		}
		defaultArgs[i] = item
	}

	port, err := strconv.Atoi(defaultArgs[2])
	if err != nil {
		port = 8080
	}

	return &Config{
		NodeName:      defaultArgs[0],
		ListenAddress: defaultArgs[1],
		ListenPort:    uint16(port),
	}
}
