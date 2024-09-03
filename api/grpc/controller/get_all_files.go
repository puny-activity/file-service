package controller

import (
	"context"
	"github.com/puny-activity/file-service/pkg/proto/gen/fileserviceproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Controller) GetAllFiles(ctx context.Context, request *fileserviceproto.GetAllFilesRequest) (*fileserviceproto.GetAllFilesResponse, error) {
	files, err := a.fileUseCase.GetAllFiles(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	filesResponse := make([]*fileserviceproto.GetAllFilesResponse_Item, len(files))
	for i, file := range files {
		filesResponse[i] = &fileserviceproto.GetAllFilesResponse_Item{
			Id:          file.ID.String(),
			Name:        file.Name,
			ContentType: file.ContentType.String(),
			Path:        file.Path,
			Size:        file.Size,
			Metadata:    file.Metadata,
			Md5:         file.MD5,
		}
	}

	return &fileserviceproto.GetAllFilesResponse{
		Files: filesResponse,
	}, nil
}
