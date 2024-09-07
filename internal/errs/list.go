package errs

import "fmt"

var (
	InvalidAPIVersion      = fmt.Errorf("invalid api version")
	InvalidFileIDParameter = fmt.Errorf("invalid file id parameter")

	FailedToStreamFile = fmt.Errorf("failed to stream file")

	Unexpected = fmt.Errorf("unexpected error")
)

// Unexpected
var unexpectedError = internalError{
	error: Unexpected,
	code:  "U-1",
}

var errorList = []internalError{
	// Request
	{
		error: InvalidAPIVersion,
		code:  "R-1",
	},
	{
		error: InvalidFileIDParameter,
		code:  "R-2",
	},
	// Runtime
	{
		error: FailedToStreamFile,
		code:  "T-1",
	},
}
