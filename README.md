# Сервис транскрипции

Сервис транскрипции - микросервис на Go, который принимает аудиопотоки по gRPC, выполняет speech-to-text с помощью вненего сервиса и публикует структурированные события транскрибирования в NATS для дальнейшей обработки.

**Ключевые возможности**
- Приём аудио в реальном времени по gRPC
- Транскрипция через Whisper (внешний эндпоинт)
- Публикация результатов в NATS

**Технологии**
- **RPC:** gRPC (google.golang.org/grpc)
- **Брокер сообщений:** NATS

**Переменные окружения**
- `WHISPERX_ENDPOINT` - URL эндпоинта Whisper для транскрипции
- `NATS_URL` - URL подключения к NATS
- `NATS_SUBJECT` - subject в NATS для публикации событий транскрипции
- `GRPC_ADDRESS` - адрес gRPC-сервиса для получения аудио


**Где смотреть код**
- [cmd/main.go](cmd/main.go) - точка входа и конфигурация сервисов
- [proto/audio_ingest.proto](proto/audio_ingest.proto) - определение API для приёма аудио
- [internal/presentation/listener.go](internal/presentation/listener.go) - слушатель NATS
- [internal/infrastructure/transcriber/whisper.go](internal/infrastructure/transcriber/whisper.go) - реализация транскрибера
- [internal/infrastructure/publisher/nats_publisher.go](internal/infrastructure/publisher/nats_publisher.go) - публикация в NATS
- [Dockerfile](Dockerfile) - Docker сборка
