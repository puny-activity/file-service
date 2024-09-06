package controller

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/puny-activity/file-service/pkg/proto/gen/fileserviceproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Controller) GetChanges(ctx context.Context, request *fileserviceproto.GetChangesRequest) (*fileserviceproto.GetChangesResponse, error) {
	if request.Since == nil {
		return nil, status.Error(codes.InvalidArgument, "since value is required")
	}

	since := carbon.FromStdTime(request.Since.AsTime())

	changes, err := a.fileUseCase.GetChangedFiles(ctx, since)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	createdFilesResponse := make([]*fileserviceproto.FileInfo, len(changes.Created))
	for i, createdFile := range changes.Created {
		createdFilesResponse[i] = &fileserviceproto.FileInfo{
			Id:          &fileserviceproto.FileInfoId{Id: createdFile.ID.String()},
			Name:        createdFile.Name,
			ContentType: createdFile.ContentType.String(),
			Path:        createdFile.Path.String(),
			Size:        createdFile.Size,
			Metadata:    createdFile.Metadata,
			Md5:         createdFile.MD5,
		}
	}
	updatedFilesResponse := make([]*fileserviceproto.FileInfo, len(changes.Updated))
	for i, updatedFile := range changes.Updated {
		updatedFilesResponse[i] = &fileserviceproto.FileInfo{
			Id:          &fileserviceproto.FileInfoId{Id: updatedFile.ID.String()},
			Name:        updatedFile.Name,
			ContentType: updatedFile.ContentType.String(),
			Path:        updatedFile.Path.String(),
			Size:        updatedFile.Size,
			Metadata:    updatedFile.Metadata,
			Md5:         updatedFile.MD5,
		}
	}
	deletedFileIDsResponse := make([]*fileserviceproto.FileInfoId, len(changes.Deleted))
	for i, deletedFileID := range changes.Deleted {
		deletedFileIDsResponse[i] = &fileserviceproto.FileInfoId{Id: deletedFileID.String()}
	}

	return &fileserviceproto.GetChangesResponse{
		CreatedFiles:   createdFilesResponse,
		UpdatedFiles:   updatedFilesResponse,
		DeletedFileIds: deletedFileIDsResponse,
	}, nil
}
