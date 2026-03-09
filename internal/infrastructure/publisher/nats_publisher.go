package publisher

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type AudioTranscribedPayload struct {
	AudioId         string `json:"audioId"`
	WorkSpaceId     string `json:"workSpaceId"`
	UploadUserId    string `json:"uploadUserId"`
	TranscriptionId string `json:"transcriptionId"`
}

type NatsPublisher struct {
	client  natsPublisherClient
	subject string
}

func NewNatsPublisher(url string, subject string) (*NatsPublisher, error) {
	natsConn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsPublisher{client: natsConn, subject: subject}, nil
}

func NewNatsPublisherFromConn(conn *nats.Conn, subject string) *NatsPublisher {
	return &NatsPublisher{client: conn, subject: subject}
}

type natsPublisherClient interface {
	PublishMsg(m *nats.Msg) error
}

func (n *NatsPublisher) Publish(ctx context.Context, audioId, workSpaceId, uploadUserid, transcription string) error {
	payload := AudioTranscribedPayload{
		AudioId:         audioId,
		WorkSpaceId:     workSpaceId,
		UploadUserId:    uploadUserid,
		TranscriptionId: transcription,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &nats.Msg{
		Subject: n.subject,
		Data:    data,
		Header:  nats.Header{"Content-Type": []string{"application/json"}},
	}

	err = n.client.PublishMsg(msg)
	if err != nil {
		return err
	}

	return nil
}
