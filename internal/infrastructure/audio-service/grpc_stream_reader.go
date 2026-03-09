package audio_service

import (
	"io"

	"github.com/mxdnght0/transcription-service/proto"
)

type GrpcStreamReader struct {
	stream proto.AudioChunkStream
	data   []byte
	closed bool
}

func newGrpcStreamReader(stream proto.AudioChunkStream) *GrpcStreamReader {
	return &GrpcStreamReader{
		stream: stream,
		data:   make([]byte, 0),
		closed: false,
	}
}

func (g GrpcStreamReader) Read(p []byte) (n int, err error) {
	if g.closed {
		return 0, io.ErrClosedPipe
	}

	if len(g.data) == 0 {
		chunk, err := g.stream.Recv()
		if err == io.EOF {
			g.closed = true
			return 0, io.EOF
		} else if err != nil {
			return 0, err
		}

		g.data = chunk.Content
		if len(g.data) == 0 {
			return g.Read(p)
		}
	}

	n = copy(p, g.data)
	if n < len(g.data) {
		g.data = g.data[n:]
	} else {
		g.data = nil
	}

	return n, nil
}

func (g GrpcStreamReader) Close() error {
	if g.closed {
		return io.ErrClosedPipe
	}
	g.closed = true

	return nil
}
