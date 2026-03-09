package interfaces

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, audioId, workSpaceId, uploadUserId, transcription string) error
}
