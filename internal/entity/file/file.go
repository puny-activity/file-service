package file

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/puny-activity/file-service/internal/entity/file/contenttype"
	"github.com/puny-activity/file-service/internal/entity/file/path"
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

func (e ID) String() string {
	return uuid.UUID(e).String()
}

type File struct {
	ID          *ID
	Path        path.Path
	Name        string
	ContentType contenttype.Type
	Size        int64
	Metadata    json.RawMessage
	MD5         string
}

func (e File) GenerateID() File {
	e.ID = util.ToPointer(ID(uuid.New()))
	return e
}

type Changes struct {
	Created []File
	Updated []File
	Deleted []ID
}
