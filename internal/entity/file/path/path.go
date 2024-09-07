package path

import (
	"fmt"
	"strings"
)

type Path struct {
	rootName     string
	relativePath string
}

func New(rootName string, relativePath string) Path {
	return Path{
		rootName:     rootName,
		relativePath: relativePath,
	}
}

func NewByString(path string) (*Path, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("there is no \"/\" in the path")
	}
	return &Path{
		rootName:     parts[0],
		relativePath: strings.TrimPrefix(path, fmt.Sprintf("%s/", parts[0])),
	}, nil
}

func (e Path) RootName() string {
	return e.rootName
}

func (e Path) RelativePath() string {
	return e.relativePath
}

func (e Path) String() string {
	return fmt.Sprintf("%s/%s", e.rootName, e.relativePath)
}
