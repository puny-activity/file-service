package miniostorage

import (
	"context"
	"encoding/json"
	"github.com/minio/minio-go/v7"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/file/filecontenttype"
	"github.com/puny-activity/file-service/pkg/metadatareader"
	"github.com/puny-activity/file-service/pkg/werr"
)

func (s *Storage) GetFiles(ctx context.Context) ([]file.File, error) {
	buckets, err := s.minioClient.ListBuckets(ctx)
	if err != nil {
		return nil, werr.WrapSE("failed to get buckets", err)
	}

	files := make([]file.File, 0)
	for _, bucket := range buckets {
		for objectInfo := range s.minioClient.ListObjects(ctx, bucket.Name, minio.ListObjectsOptions{
			Recursive:    true,
			WithMetadata: true,
		}) {
			if objectInfo.Err != nil {
				return nil, werr.WrapSE("failed to get object info", objectInfo.Err)
			}

			object, err := s.minioClient.GetObject(ctx, bucket.Name, objectInfo.Key, minio.GetObjectOptions{})
			if err != nil {
				return nil, werr.WrapSE("failed to get object", objectInfo.Err)
			}

			md5Sum, err := calculateMD5(object)
			if err != nil {
				return nil, werr.WrapSE("failed to calculate md5 sum of object", err)
			}

			contentTypeStr := objectInfo.UserMetadata["content-type"]
			if contentTypeStr == "" {
				contentTypeStr = "application/octet-stream"
			}
			contentType, err := filecontenttype.New(contentTypeStr)
			if err != nil {
				return nil, werr.WrapSE("failed to construct content type", err)
			}

			var metadata json.RawMessage
			if contentType.IsAudio() {
				metadataObject, err := metadatareader.GetAudioMetadata(object)
				if err != nil {
					return nil, werr.WrapSE("failed to get audio metadata", err)
				}
				metadata, err = metadataObject.JSONRawMessage()
				if err != nil {
					return nil, werr.WrapSE("failed to convert metadata to json", err)
				}
			} else if contentType.IsImage() {
				metadataObject, err := metadatareader.GetImageMetadata(object)
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
				Path:        bucket.Name + "/" + objectInfo.Key,
				Name:        nameByPath(objectInfo.Key),
				ContentType: contentType,
				Size:        objectInfo.Size,
				Metadata:    metadata,
				MD5:         md5Sum,
			})
		}
	}

	return files, nil
}
