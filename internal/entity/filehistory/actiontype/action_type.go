package actiontype

import "fmt"

type Type string

const (
	Unknown Type = "UNKNOWN"
	Create  Type = "CREATE"
	Update  Type = "UPDATE"
	Delete  Type = "DELETE"
)

var actionTypes = map[string]Type{
	"CREATE": Create,
	"UPDATE": Update,
	"DELETE": Delete,
}

func New(typeName string) (Type, error) {
	actionType, ok := actionTypes[typeName]
	if !ok {
		return Unknown, fmt.Errorf("unknown action type: %s", typeName)
	}
	return actionType, nil
}

func (e Type) String() string {
	return string(e)
}
