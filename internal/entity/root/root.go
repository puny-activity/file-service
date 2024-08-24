package root

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/puny-activity/file-service/internal/entity/root/roottype"
	"github.com/puny-activity/file-service/pkg/util"
)

type ID uuid.UUID

func (e ID) String() string {
	return uuid.UUID(e).String()
}

type Root struct {
	ID      *ID
	Type    roottype.Type
	Name    string
	Details json.RawMessage
}

func (e Root) GenerateID() Root {
	e.ID = util.ToPointer(ID(uuid.New()))
	return e
}

func (e Root) MarshalJSON() ([]byte, error) {
	return nil, nil
}
