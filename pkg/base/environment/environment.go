package environment

import (
	"errors"
	"strings"
)

type Environment string

const (
	Production  Environment = "production"
	Test        Environment = "test"
	Development Environment = "development"
	Local       Environment = "local"
	Unknown     Environment = "unknown"
)

func New(name string) (Environment, error) {
	switch strings.ToLower(name) {
	case "production", "prod", "p":
		return Production, nil
	case "test", "t":
		return Test, nil
	case "development", "dev", "d":
		return Development, nil
	case "local", "l":
		return Local, nil
	default:
		return Unknown, errors.New("failed to parse environment")
	}
}

func (e Environment) String() string {
	return string(e)
}

func (e Environment) IsProduction() bool {
	return e.String() == Production.String()
}

func (e Environment) IsTest() bool {
	return e.String() == Test.String()
}

func (e Environment) IsDevelopment() bool {
	return e.String() == Development.String()
}

func (e Environment) IsLocal() bool {
	return e.String() == Local.String()
}
