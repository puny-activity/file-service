package file

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/puny-activity/file-service/internal/entity/file/filecontenttype"
	"github.com/puny-activity/file-service/pkg/util"
)

type ID uuid.UUID

func (e ID) String() string {
	return uuid.UUID(e).String()
}

type File struct {
	ID          *ID
	Path        string
	Name        string
	ContentType filecontenttype.Type
	Size        int64
	Metadata    json.RawMessage
	MD5         string
}

func (e File) GenerateID() File {
	e.ID = util.ToPointer(ID(uuid.New()))
	return e
}
