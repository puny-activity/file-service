package localstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/file/filecontenttype"
	"github.com/puny-activity/file-service/pkg/metadatareader"
	"github.com/puny-activity/file-service/pkg/werr"
	"os"
	"path/filepath"
	"strings"
)

func (s *Storage) GetFiles(ctx context.Context) ([]file.File, error) {
	files, err := getFilesByPath(s.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get files by path: %w", err)
	}

	for i := range files {
		files[i].Path = strings.TrimPrefix(files[i].Path, s.basePath)
	}

	return files, nil
}

func getFilesByPath(path string) ([]file.File, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	files := make([]file.File, 0)

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			subDirFiles, err := getFilesByPath(entryPath)
			if err != nil {
				return nil, fmt.Errorf("failed to get subdir files by path: %w", err)
			}
			files = append(files, subDirFiles...)
		} else {
			fileEntry, err := os.Open(entryPath)
			if err != nil {
				return nil, fmt.Errorf("failed to open file: %w", err)
			}

			fileInfo, err := entry.Info()
			if err != nil {
				return nil, fmt.Errorf("failed to stat file: %w", err)
			}

			md5Sum, err := calculateMD5(entryPath)
			if err != nil {
				return nil, werr.WrapSE("failed to calculate md5 sum of object", err)
			}

			contentTypeStr, err := detectContentType(fileEntry)
			if err != nil {
				return nil, werr.WrapSE("failed to detect content type", err)
			}
			contentType, err := filecontenttype.New(contentTypeStr)
			if err != nil {
				return nil, werr.WrapSE("failed to construct content type", err)
			}

			var metadata json.RawMessage
			if contentType.IsAudio() {
				metadataObject, err := metadatareader.GetAudioMetadata(fileEntry)
				if err != nil {
					return nil, werr.WrapSE("failed to get audio metadata", err)
				}
				metadata, err = metadataObject.JSONRawMessage()
				if err != nil {
					return nil, werr.WrapSE("failed to convert metadata to json", err)
				}
			} else if contentType.IsImage() {
				metadataObject, err := metadatareader.GetImageMetadata(fileEntry)
				if err != nil {
					return nil, werr.WrapSE("failed to get image metadata", err)
				}
				metadata, err = metadataObject.JSONRawMessage()
				if err != nil {
					return nil, werr.WrapSE("failed to convert metadata to json", err)
				}
			} else {
				metadata = []byte("{}")
			}

			files = append(files, file.File{
				Path:        entryPath,
				Name:        entry.Name(),
				ContentType: contentType,
				Size:        fileInfo.Size(),
				Metadata:    metadata,
				MD5:         md5Sum,
			})
		}
	}

	return files, nil
}
