package werr

import "fmt"

func WrapES(highError error, lowString string) error {
	return fmt.Errorf("%w: %s", highError, lowString)
}

func WrapSE(highString string, lowError error) error {
	return fmt.Errorf("%s: %w", highString, lowError)
}
