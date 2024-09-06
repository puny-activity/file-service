package path

import (
	"fmt"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/pkg/werr"
	"strings"
)

type Path struct {
	rootID       root.ID
	relativePath string
}

func New(root root.ID, relativePath string) Path {
	return Path{
		rootID:       root,
		relativePath: relativePath,
	}
}

func NewByString(path string) (*Path, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("there is no \"/\" in the path")
	}
	rootID, err := root.ParseID(parts[0])
	if err != nil {
		return nil, werr.WrapSE("failed to construct root id", err)
	}
	return &Path{
		rootID:       rootID,
		relativePath: strings.TrimPrefix(path, fmt.Sprintf("%s/", parts[0])),
	}, nil
}

func (e Path) RootID() root.ID {
	return e.rootID
}

func (e Path) RelativePath() string {
	return e.relativePath
}

func (e Path) String() string {
	return fmt.Sprintf("%s/%s", e.rootID.String(), e.relativePath)
}
