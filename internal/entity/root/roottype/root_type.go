package roottype

import "fmt"

type Type string

const (
	Unknown Type = "UNKNOWN"
	Local   Type = "LOCAL"
	Minio   Type = "MINIO"
)

var rootTypes = map[string]Type{
	"LOCAL": Local,
	"MINIO": Minio,
}

func Parse(typeName string) (Type, error) {
	rootType, ok := rootTypes[typeName]
	if !ok {
		return Unknown, fmt.Errorf("unknown root type: %s", typeName)
	}
	return rootType, nil
}

func (e Type) String() string {
	return string(e)
}
