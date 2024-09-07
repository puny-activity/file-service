package httpcontroller

import (
	"github.com/go-chi/chi/v5"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/errs"
	"io"
	"net/http"
)

func (c *Controller) StreamFile(w http.ResponseWriter, r *http.Request) error {
	fileIDStr := chi.URLParam(r, "file_id")
	if fileIDStr == "" {
		return errs.InvalidFileIDParameter
	}
	fileID, err := file.ParseID(fileIDStr)
	if err != nil {
		return errs.InvalidFileIDParameter
	}

	fileStream, fileInfo, err := c.app.FileUseCase.StreamFile(r.Context(), fileID)
	if err != nil {
		return err
	}
	defer func(fileStream io.ReadCloser) {
		err := fileStream.Close()
		if err != nil {
			c.log.Error().Err(err).Str("fileId", fileID.String()).Msg("failed to close file")
		}
	}(fileStream)

	w.Header().Set("Content-Type", fileInfo.ContentType.String())
	w.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name)
	w.WriteHeader(http.StatusOK)

	_, err = io.Copy(w, fileStream)
	if err != nil {
		return errs.FailedToStreamFile
	}

	return nil
}
