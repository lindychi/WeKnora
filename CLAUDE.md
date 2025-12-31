# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WeKnora is an LLM-powered RAG (Retrieval-Augmented Generation) framework for document understanding and semantic retrieval. It handles complex, heterogeneous documents through multimodal preprocessing, semantic vector indexing, intelligent retrieval, and LLM inference.

## Development Commands

### Quick Development Mode (Recommended)

```bash
# Start infrastructure services (PostgreSQL, Redis, MinIO, DocReader, etc.)
make dev-start

# Start backend (separate terminal) - supports Air hot-reload
make dev-app

# Start frontend (separate terminal) - Vite with hot-reload
make dev-frontend

# Check service status
make dev-status

# Stop all services
make dev-stop
```

### Full Docker Deployment

```bash
# Start all services (includes Ollama)
./scripts/start_all.sh   # or: make start-all

# Stop all services
./scripts/start_all.sh --stop   # or: make stop-all

# Optional profiles
docker-compose --profile full up -d      # All features
docker-compose --profile neo4j up -d     # With knowledge graph
docker-compose --profile minio up -d     # With MinIO storage
```

### Build & Test

```bash
# Build Go application
make build

# Run tests
make test   # or: go test -v ./...

# Lint code
make lint   # uses golangci-lint

# Format code
make fmt

# Generate Swagger docs
make docs
```

### Database Migrations

```bash
make migrate-up                    # Apply migrations
make migrate-down                  # Rollback
make migrate-create name=<name>    # Create new migration
make migrate-version              # Check current version
```

### Docker Image Build

```bash
make docker-build-app         # Backend image
make docker-build-docreader   # Document reader image
make docker-build-frontend    # Frontend image
make docker-build-all         # All images
```

## Architecture

### System Layers

```
┌──────────────────────────────────────────────────────────────┐
│                        Frontend (Vue 3)                       │
│    TDesign UI · Vue Router · Pinia · vue-i18n · Vite         │
├──────────────────────────────────────────────────────────────┤
│                     Backend API (Go/Gin)                      │
│  cmd/server/main.go → router → handler → service → repo      │
├──────────────────────────────────────────────────────────────┤
│                    DocReader (Python/gRPC)                    │
│     PDF/Word/Markdown parsing · OCR · URL extraction          │
├──────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                       │
│  PostgreSQL(pgvector) · Redis · MinIO/COS · Neo4j · Ollama   │
└──────────────────────────────────────────────────────────────┘
```

### Backend Structure (`internal/`)

- **container/**: Dependency injection (uber/dig) - wires all components
- **handler/**: HTTP request handlers (Gin framework)
- **application/service/**: Business logic layer
- **application/repository/**: Data access layer with retriever implementations
- **types/**: Domain types, interfaces, and DTOs
- **models/**: LLM provider integrations (embedding, chat, rerank)
- **agent/**: ReACT agent implementation with tool calling
- **mcp/**: MCP (Model Context Protocol) client manager
- **event/**: Event bus for async processing
- **stream/**: SSE streaming for chat responses

### Key Patterns

1. **Chat Pipeline** (`internal/application/service/chat_pipline/`): Plugin-based architecture for chat processing (search, rerank, rewrite, completion, streaming)

2. **Retriever Registry** (`internal/application/service/retriever/`): Supports multiple vector stores:
   - PostgreSQL with pgvector
   - Elasticsearch v7/v8
   - Qdrant

3. **File Storage** (`internal/application/service/file/`): Abstracted file service supporting MinIO, Tencent COS, and local filesystem

4. **Async Tasks**: Redis-based queue with hibiken/asynq for background processing

### Frontend Structure (`frontend/src/`)

- **views/**: Page components organized by feature
- **components/**: Reusable UI components
- **api/**: Backend API client modules
- **stores/**: Pinia state management
- **i18n/**: Internationalization (en-US, zh-CN, ko-KR)

### DocReader Structure (`docreader/`)

Python gRPC service for document parsing:
- **parser/**: Document format parsers (PDF, Word, Markdown, HTML, images)
- **splitter/**: Text chunking logic
- **proto/**: gRPC service definitions

## Configuration

Environment variables in `.env`:
- `DB_DRIVER`: postgres (primary database)
- `RETRIEVE_DRIVER`: postgres, elasticsearch_v7, elasticsearch_v8, qdrant
- `STORAGE_TYPE`: local, minio, cos
- `OLLAMA_BASE_URL`: LLM inference endpoint
- `ENABLE_GRAPH_RAG`: Enable Neo4j knowledge graph

Config file: `config/config.yaml` - conversation prompts, retrieval thresholds, web search settings

## API Documentation

- Swagger UI: `http://localhost:8080/swagger/index.html` (debug mode only)
- API Base: `/api/v1`
- Auth: Bearer token (JWT) or X-API-Key header

## Key Dependencies

**Backend (Go 1.24+)**:
- Gin (HTTP framework)
- GORM (ORM with PostgreSQL/pgvector)
- uber/dig (dependency injection)
- hibiken/asynq (async task queue)
- sashabaranov/go-openai (LLM client)
- mark3labs/mcp-go (MCP protocol)

**Frontend**:
- Vue 3 + TypeScript
- TDesign Vue Next (UI library)
- Vite 7

**DocReader (Python 3.10+)**:
- Uses `uv` for package management
- PaddleOCR for image text extraction
- pdfplumber, python-docx for document parsing
