package metadatareader

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type ImageMetadata struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (m *ImageMetadata) JSONRawMessage() (json.RawMessage, error) {
	return json.Marshal(m)
}

func GetImageMetadata(file io.ReadSeeker) (ImageMetadata, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return ImageMetadata{}, fmt.Errorf("failed to seek file: %w", err)
	}

	var img image.Image
	img, _, err = image.Decode(file)
	if err != nil {
		return ImageMetadata{}, fmt.Errorf("failed to decode image: %v", err)
	}

	return ImageMetadata{
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
	}, nil
}
