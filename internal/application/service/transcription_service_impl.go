package service

import (
	"context"

	"github.com/transcription-service/internal/application/interfaces"
)

type transcriptionService struct {
	transcriber  interfaces.Transcriber
	publisher    interfaces.Publisher
	audioService interfaces.AudioService
}

func NewTranscriptionService(transcriber interfaces.Transcriber,
	publisher interfaces.Publisher,
	audioService interfaces.AudioService) interfaces.TranscriptionService {
	return &transcriptionService{
		transcriber:  transcriber,
		publisher:    publisher,
		audioService: audioService,
	}
}

func (t *transcriptionService) Transcript(ctx context.Context, audioId, workspaceId, uploadUserId string) error {
	audio, err := t.audioService.GetAudio(audioId, workspaceId, uploadUserId)
	if err != nil {
		return err
	}

	transcription, err := t.transcriber.Transcript(ctx, audio)
	if err != nil {
		return err
	}

	err = t.publisher.Publish(ctx, audioId, workspaceId, uploadUserId, transcription)
	if err != nil {
		return err
	}

	return nil
}
