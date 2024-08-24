package miniostorage

import (
	"crypto/md5"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/puny-activity/file-service/pkg/werr"
	"io"
	"strings"
)

func calculateMD5(minioObject *minio.Object) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, minioObject); err != nil {
		return "", werr.WrapSE("while calculating SHA256", err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func nameByPath(path string) string {
	pathParts := strings.Split(path, "/")
	return pathParts[len(pathParts)-1]
}
