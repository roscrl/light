package config

import (
	"fmt"
)

type Environment int

const (
	LOCAL Environment = iota
	DEV
	PROD
)

func (e Environment) String() string {
	switch e {
	case LOCAL:
		return "LOCAL"
	case DEV:
		return "DEV"
	case PROD:
		return "PROD"
	default:
		return "UNKNOWN"
	}
}

func parseEnvironment(value string) (Environment, error) {
	switch value {
	case "LOCAL":
		return LOCAL, nil
	case "DEV":
		return DEV, nil
	case "PROD":
		return PROD, nil
	default:
		return 0, fmt.Errorf("unknown environment: %s", value)
	}
}
