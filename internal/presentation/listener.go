package presentation

import (
	"context"
	"encoding/json"
	"log"

	"github.com/transcription-service/internal/application/interfaces"
	"github.com/nats-io/nats.go"
)

type AudioUploadedPayload struct {
	AudioId         string `json:"audioId"`
	WorkSpaceId     string `json:"workSpaceId"`
	TranscriptionId string `json:"transcriptionId"`
}

type Listener struct {
	natsConn             *nats.Conn
	ch                   chan *nats.Msg
	sub                  *nats.Subscription
	transcriptionService interfaces.TranscriptionService
}

func NewListener(natsURL, subject string, transcriptionService interfaces.TranscriptionService) (*Listener, error) {
	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	ch := make(chan *nats.Msg, 64)
	sub, err := natsConn.ChanSubscribe(subject, ch)
	if err != nil {
		return nil, err
	}

	return &Listener{natsConn: natsConn, ch: ch, sub: sub, transcriptionService: transcriptionService}, nil
}

func (l *Listener) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("listener stopped:", ctx.Err())
			_ = l.sub.Unsubscribe()
			return
		case msg := <-l.ch:
			go func(msg *nats.Msg) {
				audioUploaded := &AudioUploadedPayload{}
				if err := json.Unmarshal(msg.Data, audioUploaded); err != nil {
					log.Println("error unmarshalling audio uploaded payload:", err)
				}
				err := l.transcriptionService.Transcript(
					ctx,
					audioUploaded.AudioId,
					audioUploaded.WorkSpaceId,
					audioUploaded.TranscriptionId,
				)
				if err != nil {
					log.Println("error transcribing audio uploaded:", err)
				}
			}(msg)
		}
	}
}
