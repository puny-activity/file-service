package wsserver

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/puny-activity/file-service/config"
	"github.com/puny-activity/file-service/internal/app"
	"github.com/puny-activity/file-service/internal/entity/file"
	"io"
	"log"
	"net/http"
)

type WebSocketServer struct {
	cfg         *config.WebSocket
	application *app.App
}

func New(cfg *config.WebSocket, application *app.App) *WebSocketServer {
	return &WebSocketServer{
		cfg:         cfg,
		application: application,
	}
}

func (s *WebSocketServer) Start() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  s.cfg.ReadBufferSize,
			WriteBufferSize: s.cfg.WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		defer conn.Close()

		fileUUID, _ := uuid.Parse("a3a2137c-ac08-405a-8445-0f64f4bdf343")
		fileReader, err := s.application.FileUseCase.ReadFile(r.Context(), file.ID(fileUUID))
		if err != nil {
			log.Println(err)
		}
		defer fileReader.Close()

		buffer := make([]byte, s.cfg.ReadBufferSize)
		for {
			n, err := fileReader.Read(buffer)
			if err != nil && err != io.EOF {
				log.Println(err)
				break
			}
			if n == 0 {
				break
			}
			err = conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
			if err != nil {
				log.Println(err)
				break
			}
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.Port), nil)
	if err != nil {
		log.Println(err)
	}
}
