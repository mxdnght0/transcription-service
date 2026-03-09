package transcriber

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Whisper struct {
	endpoint string
	client   *http.Client
}

func NewWhisper(endpoint string) *Whisper {
	return &Whisper{
		endpoint: strings.TrimRight(endpoint, "/"),
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (w *Whisper) Transcript(ctx context.Context, audio io.ReadCloser) (string, error) {
	defer audio.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.endpoint, audio)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := w.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("transcribe failed: status %d: %s", resp.StatusCode, string(body))
	}

	var out struct {
		Success  bool   `json:"success"`
		Language string `json:"language"`
		Segments []struct {
			Start float64 `json:"start"`
			End   float64 `json:"end"`
			Text  string  `json:"text"`
		} `json:"segments"`
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&out); err != nil {
		return "", err
	}
	if !out.Success {
		return "", fmt.Errorf("transcription service returned success=false")
	}

	var b strings.Builder
	for i, s := range out.Segments {
		if i > 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "%.2f-%.2f: %s", s.Start, s.End, strings.TrimSpace(s.Text))
	}

	return b.String(), nil
}
