package httpresp

import (
	"bytes"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"time"
)

type ErrorStorage interface {
	HTTPStatusCode(err error) int
	LocalizationCode(err error) string
	Message(lang string, err error) string
}

type Wrapper struct {
	writer       *Writer
	errorStorage ErrorStorage
	log          *zerolog.Logger
}

func NewWrapper(writer *Writer, errorStorage ErrorStorage, log *zerolog.Logger) *Wrapper {
	return &Wrapper{
		writer:       writer,
		errorStorage: errorStorage,
		log:          log,
	}
}

func (w *Wrapper) Wrap(controllerFunction func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

		var requestBody bytes.Buffer
		if request.Body != nil {
			_, err := io.Copy(&requestBody, request.Body)
			if err != nil {
				w.log.Error().Err(err).Msg("failed to read request body")
				return
			}
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.log.Error().Err(err).Msg("failed to close request body")
			}
		}(request.Body)

		clonedRequest := request.Clone(request.Context())
		clonedRequest.Body = io.NopCloser(bytes.NewReader(requestBody.Bytes()))

		ctx := request.Context()
		clonedRequest = clonedRequest.WithContext(ctx)

		controllerError := controllerFunction(ww, clonedRequest)
		if controllerError != nil {
			err := w.writer.Write(
				ww,
				w.errorStorage.HTTPStatusCode(controllerError),
				ErrorResponse{
					Code:    w.errorStorage.LocalizationCode(controllerError),
					Message: controllerError.Error(),
					//Message: w.errorStorage.Message(lang, controllerError),
				})
			if err != nil {
				w.log.Debug().Err(err).Msg("failed to write error response")
			}
		}

		duration := time.Since(start)
		var requestBodyJSON []byte
		if requestBody.Len() > 0 {
			requestBodyJSON = requestBody.Bytes()
		} else {
			requestBodyJSON = []byte("null")
		}

		if ww.Status() >= http.StatusInternalServerError {
			w.log.Error().
				Err(controllerError).
				Str("method", request.Method).
				Str("path", request.URL.String()).
				Str("duration", duration.String()).
				RawJSON("request_body", requestBodyJSON).
				Int("request_body_length_bytes", requestBody.Len()).
				Int("response_body_length_bytes", ww.BytesWritten()).
				Int("response_status", ww.Status()).
				Str("user_agent", request.UserAgent()).
				Str("source_ip", request.RemoteAddr).
				Msgf("request handled with unexpected error: %d %s %s", ww.Status(), request.Method, request.URL.Path)
		} else if ww.Status() >= http.StatusBadRequest {
			w.log.Warn().
				Err(controllerError).
				Str("method", request.Method).
				Str("path", request.URL.String()).
				Str("duration", duration.String()).
				RawJSON("request_body", requestBodyJSON).
				Int("request_body_length_bytes", requestBody.Len()).
				Int("response_body_length_bytes", ww.BytesWritten()).
				Int("response_status", ww.Status()).
				Str("user_agent", request.UserAgent()).
				Str("source_ip", request.RemoteAddr).
				Msgf("request handled with error: %d %s %s", ww.Status(), request.Method, request.URL.Path)
		} else {
			w.log.Info().
				Str("method", request.Method).
				Str("path", request.URL.String()).
				Str("duration", duration.String()).
				RawJSON("request_body", requestBodyJSON).
				Int("request_body_length_bytes", requestBody.Len()).
				Int("response_body_length_bytes", ww.BytesWritten()).
				Int("response_status", ww.Status()).
				Str("user_agent", request.UserAgent()).
				Str("source_ip", request.RemoteAddr).
				Msgf("request handled succussfully: %d %s %s", ww.Status(), request.Method, request.URL.Path)
		}
	}
}

func (w *Wrapper) WrapAnonymous(controllerFunction func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

		var requestBody bytes.Buffer
		if request.Body != nil {
			_, err := io.Copy(&requestBody, request.Body)
			if err != nil {
				w.log.Error().Err(err).Msg("failed to read request body")
				return
			}
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.log.Error().Err(err).Msg("failed to close request body")
			}
		}(request.Body)

		clonedRequest := request.Clone(request.Context())
		clonedRequest.Body = io.NopCloser(bytes.NewReader(requestBody.Bytes()))

		//lang := request.Context().Value(ctxbase.Language).(string)

		ctx := request.Context()
		clonedRequest = clonedRequest.WithContext(ctx)

		controllerError := controllerFunction(ww, clonedRequest)
		if controllerError != nil {
			err := w.writer.Write(
				ww,
				w.errorStorage.HTTPStatusCode(controllerError),
				ErrorResponse{
					Code:    w.errorStorage.LocalizationCode(controllerError),
					Message: controllerError.Error(),
					//Message: w.errorStorage.Message(lang, controllerError),
				})
			if err != nil {
				w.log.Debug().Err(err).Msg("failed to write error response")
			}
		}

		duration := time.Since(start)

		if ww.Status() >= http.StatusInternalServerError {
			w.log.Error().
				Err(controllerError).
				Str("method", request.Method).
				Str("path", request.URL.String()).
				Str("duration", duration.String()).
				Int("request_body_length_bytes", requestBody.Len()).
				Int("response_body_length_bytes", ww.BytesWritten()).
				Int("response_status", ww.Status()).
				Str("user_agent", request.UserAgent()).
				Str("source_ip", request.RemoteAddr).
				Msgf("request handled with unexpected error: %d %s %s", ww.Status(), request.Method, request.URL.Path)
		} else if ww.Status() >= http.StatusBadRequest {
			w.log.Warn().
				Err(controllerError).
				Str("method", request.Method).
				Str("path", request.URL.String()).
				Str("duration", duration.String()).
				Int("request_body_length_bytes", requestBody.Len()).
				Int("response_body_length_bytes", ww.BytesWritten()).
				Int("response_status", ww.Status()).
				Str("user_agent", request.UserAgent()).
				Str("source_ip", request.RemoteAddr).
				Msgf("request handled with error: %d %s %s", ww.Status(), request.Method, request.URL.Path)
		} else {
			w.log.Info().
				Str("method", request.Method).
				Str("path", request.URL.String()).
				Str("duration", duration.String()).
				Int("request_body_length_bytes", requestBody.Len()).
				Int("response_body_length_bytes", ww.BytesWritten()).
				Int("response_status", ww.Status()).
				Str("user_agent", request.UserAgent()).
				Str("source_ip", request.RemoteAddr).
				Msgf("request handled succussfully: %d %s %s", ww.Status(), request.Method, request.URL.Path)
		}
	}
}

func (w *Wrapper) WrapWithoutLog(controllerFunction func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ww := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

		//lang := request.Context().Value(ctxbase.Language).(string)

		controllerError := controllerFunction(ww, request)
		if controllerError != nil {
			err := w.writer.Write(
				ww,
				w.errorStorage.HTTPStatusCode(controllerError),
				ErrorResponse{
					Code:    w.errorStorage.LocalizationCode(controllerError),
					Message: controllerError.Error(),
					//Message: w.errorStorage.Message(lang, controllerError),
				})
			if err != nil {
				w.log.Debug().Err(err).Msg("failed to write error response")
			}
		}
	}
}
