package audio_service

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/transcription-service/proto"
)

type fakeStream struct {
	chunks [][]byte
	i      int
}

func (f *fakeStream) Recv() (*proto.AudioChunk, error) {
	if f.i >= len(f.chunks) {
		return nil, io.EOF
	}
	c := f.chunks[f.i]
	f.i++
	return &proto.AudioChunk{Content: c}, nil
}

type fakeClient struct {
	stream proto.AudioChunkStream
	err    error
}

func (f *fakeClient) GetAudio(ctx context.Context, req *proto.GetAudioRequest) (proto.AudioChunkStream, error) {
	return f.stream, f.err
}

func TestGetAudio_ReadsChunks(t *testing.T) {
	fs := &fakeStream{chunks: [][]byte{[]byte("ab"), []byte("cd"), []byte("")}}
	fc := &fakeClient{stream: fs}
	svc := NewGrpcAudioService(fc)

	r, err := svc.GetAudio("a", "w", "u")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer r.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil && err != io.EOF {
		t.Fatalf("copy error: %v", err)
	}
	if buf.String() != "abcd" {
		t.Fatalf("unexpected data: %q", buf.String())
	}
}

func TestGetAudio_GetAudioError(t *testing.T) {
	fc := &fakeClient{err: io.ErrUnexpectedEOF}
	svc := NewGrpcAudioService(fc)
	_, err := svc.GetAudio("a", "w", "u")
	if err == nil {
		t.Fatal("expected error")
	}
}
