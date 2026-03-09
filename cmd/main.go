package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/transcription-service/internal/application/service"
	audio_service "github.com/transcription-service/internal/infrastructure/audio-service"
	"github.com/transcription-service/internal/infrastructure/publisher"
	"github.com/transcription-service/internal/infrastructure/transcriber"
	"github.com/transcription-service/internal/presentation"
	"github.com/transcription-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("Starting transcription service")

	endpoint := os.Getenv("WHISPERX_ENDPOINT")
	if endpoint == "" {
		log.Fatal("Whisper endpoint not set")
	}
	t := transcriber.NewWhisper(endpoint)

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL not set")
	}
	subject := os.Getenv("NATS_SUBJECT")
	if subject == "" {
		log.Fatal("NATS_SUBJECT not set")
	}
	p, err := publisher.NewNatsPublisher(natsURL, subject)
	if err != nil {
		log.Fatal("failed to create nats publisher: ", err)
	}

	grpcAddress := os.Getenv("GRPC_ADDRESS")
	if grpcAddress == "" {
		log.Fatal("GRPC_ADDRESS not set")
	}
	conn, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to create grpc connection: ", err)
	}
	defer conn.Close()
	client := proto.NewAudioIngestServiceClient(conn)
	audioService := audio_service.NewGrpcAudioService(client)

	transcriptionService := service.NewTranscriptionService(t, p, audioService)

	listener, err := presentation.NewListener(natsURL, subject, transcriptionService)
	if err != nil {
		log.Fatal("failed to create listener: ", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		listener.Listen(ctx)
	}()

	<-ctx.Done()

	log.Println("Shutting down")
}
