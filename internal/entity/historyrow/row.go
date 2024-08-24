package historyrow

import (
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/historyrow/actiontype"
	"github.com/puny-activity/file-service/pkg/util"
)

type ID uuid.UUID

func (e ID) String() string {
	return uuid.UUID(e).String()
}

type Row struct {
	ID          *ID
	FileID      file.ID
	Action      actiontype.Type
	PerformedAt carbon.Carbon
}

func (e Row) GenerateID() Row {
	e.ID = util.ToPointer(ID(uuid.New()))
	return e
}
