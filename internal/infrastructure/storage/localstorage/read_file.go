package localstorage

import (
	"context"
	"github.com/puny-activity/file-service/internal/entity/file"
	"io"
	"os"
	"path/filepath"
)

func (s *Storage) ReadFile(ctx context.Context, file file.File) (io.ReadCloser, error) {
	return os.Open(filepath.Join(s.basePath, file.Path.RelativePath()))
}
