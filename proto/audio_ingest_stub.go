package proto

import "context"


type GetAudioRequest struct {
	AudioId        string
	WorkspaceId    string
	UploaderUserId string
}

type AudioChunk struct {
	AudioId     string
	WorkspaceId string
	Content     []byte
}

type AudioChunkStream interface {
	Recv() (*AudioChunk, error)
}

type AudioIngestServiceClient interface {
	GetAudio(ctx context.Context, req *GetAudioRequest) (AudioChunkStream, error)
}

func NewAudioIngestServiceClient(_ interface{}) AudioIngestServiceClient {
	return nil
}
