package root

import (
	"encoding/json"
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"github.com/puny-activity/file-service/internal/entity/root/roottype"
	"github.com/puny-activity/file-service/pkg/util"
	"github.com/puny-activity/file-service/pkg/werr"
)

type ID uuid.UUID

func ParseID(id string) (ID, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return ID{}, werr.WrapSE("failed to parse uuid", err)
	}
	return ID(idUUID), nil
}

func MustNewID(id string) ID {
	return ID(uuid.MustParse(id))
}

func (e ID) String() string {
	return uuid.UUID(e).String()
}

type Root struct {
	ID        *ID
	Type      roottype.Type
	Name      string
	Config    json.RawMessage
	ScannedAt *carbon.Carbon
}

func (e Root) GenerateID() Root {
	e.ID = util.ToPointer(ID(uuid.New()))
	return e
}
