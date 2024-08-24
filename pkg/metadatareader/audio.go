package metadatareader

import (
	"encoding/json"
	"fmt"
	"github.com/puny-activity/file-service/pkg/util"
	"github.com/wtolson/go-taglib"
	"io"
	"os"
)

type AudioMetadata struct {
	Title        *string `json:"title,omitempty"`
	DurationNs   int64   `json:"durationNs,omitempty"`
	Artist       *string `json:"artist,omitempty"`
	Album        *string `json:"album,omitempty"`
	Genre        *string `json:"genre,omitempty"`
	Year         *int    `json:"year,omitempty"`
	TrackNumber  *int    `json:"trackNumber,omitempty"`
	Comment      *string `json:"comment,omitempty"`
	Channels     int     `json:"channels,omitempty"`
	BitrateKbps  int     `json:"bitrateKbps,omitempty"`
	SampleRateHz int     `json:"sampleRateHz,omitempty"`
}

func (m *AudioMetadata) JSONRawMessage() (json.RawMessage, error) {
	return json.Marshal(m)
}

func GetAudioMetadata(file io.ReadSeeker) (AudioMetadata, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return AudioMetadata{}, fmt.Errorf("failed to seek audio metadata file: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "tmp-audio-*")
	if err != nil {
		return AudioMetadata{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, file); err != nil {
		tmpFile.Close()
		return AudioMetadata{}, fmt.Errorf("failed to copy audio metadata to temp file: %w", err)
	}
	tmpFile.Close()

	tagFile, err := taglib.Read(tmpFile.Name())
	if err != nil {
		return AudioMetadata{}, fmt.Errorf("failed to read tag metadata from temp file: %w", err)
	}
	defer tagFile.Close()

	var title *string = nil
	if tagFile.Title() != "" {
		title = util.ToPointer(tagFile.Title())
	}

	var artist *string = nil
	if tagFile.Artist() != "" {
		artist = util.ToPointer(tagFile.Artist())
	}

	var album *string = nil
	if tagFile.Album() != "" {
		album = util.ToPointer(tagFile.Album())
	}

	var genre *string = nil
	if tagFile.Genre() != "" {
		genre = util.ToPointer(tagFile.Genre())
	}

	var year *int = nil
	if tagFile.Year() != 0 {
		year = util.ToPointer(tagFile.Year())
	}

	var trackNumber *int = nil
	if tagFile.Track() != 0 {
		trackNumber = util.ToPointer(tagFile.Track())
	}

	durationNs := tagFile.Length().Nanoseconds()

	sampleRateHz := tagFile.Samplerate()

	channels := tagFile.Channels()

	bitrateKbps := tagFile.Bitrate()

	var comment *string = nil
	if tagFile.Comment() != "" {
		comment = util.ToPointer(tagFile.Comment())
	}

	return AudioMetadata{
		Title:        title,
		DurationNs:   durationNs,
		Artist:       artist,
		Album:        album,
		Genre:        genre,
		Year:         year,
		TrackNumber:  trackNumber,
		Comment:      comment,
		Channels:     channels,
		BitrateKbps:  bitrateKbps,
		SampleRateHz: sampleRateHz,
	}, nil
}
