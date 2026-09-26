# Go Document API

A Go backend for an app where you can upload documents and ask questions about their contents.

The goal is simple: upload a long document, ask a question, and get an answer based on the relevant parts of that document. This is a learning project and is still being built.

## What works today

- Start an HTTP server and check it through `/health`.
- Connect to PostgreSQL when the server starts.
- Upload a file and save it in the local `uploads/` folder.
- Store the file's name, storage path, status, and creation time in PostgreSQL.
- Return the saved document's details as JSON.

The upload request has a 10 MiB limit, including the form data around the file. Files get unique storage names so two uploads with the same original filename do not overwrite each other.

There is also code to split text into overlapping chunks, save those chunks, and update document status. These pieces are not connected to the upload flow yet, so uploaded documents stay `pending`.

AI answers, embeddings, document search, and automatic text processing are not implemented yet. There is no frontend or authentication yet.

## How it works

```text
Upload request
    -> HTTP handler receives the file
    -> Service saves the file in uploads/
    -> Repository saves its details in PostgreSQL
    -> API returns document details
```

The file itself lives on disk. The database stores information about the file.

## Tools used

| Tool | What it does |
| --- | --- |
| Go | Runs the server and application logic |
| PostgreSQL | Stores document details and, later, text chunks |
| Docker Compose | Starts the local database with its settings |
| pgxpool | Lets Go communicate with PostgreSQL using reusable connections |

## Run locally

You need Go 1.25 or newer and Docker with Docker Compose. Run the commands below from the project folder. The Go server runs on your computer; PostgreSQL runs in Docker.

### 1. Configure the database

If you do not already have a `.env` file, copy the example:

```bash
cp .env.example .env
```

Edit `.env` and replace the placeholder with your local database password. If you already initialized the database, use its existing password; changing this file does not change the stored database password.

`.env` is ignored by Git. `.env.example` contains only a placeholder and can be shared.

### 2. Start PostgreSQL

```bash
docker compose up -d
docker compose exec postgres pg_isready -U document -d document_ai
```

Wait for `accepting connections`. The database is available to your Go app at `127.0.0.1:5433`. Docker stores its data in a named volume.

### 3. Create the tables

Run each migration once, in this order. Skip any migration you have already applied.

```bash
docker compose exec -T postgres psql -U document -d document_ai \
  -v ON_ERROR_STOP=1 < migrations/001_create_documents.sql

docker compose exec -T postgres psql -U document -d document_ai \
  -v ON_ERROR_STOP=1 < migrations/002_create_chunks.sql
```

### 4. Start the API

Replace `YOUR_PASSWORD` below with the same password used in `.env`. URL-encode special characters in the password when putting it in this connection URL.

```bash
export DATABASE_URL='postgres://document:YOUR_PASSWORD@127.0.0.1:5433/document_ai?sslmode=disable'
go run ./cmd/api
```

`sslmode=disable` is for this local setup. Docker Compose reads `.env`, but the Go application does not load it automatically. Set `DATABASE_URL` in the same terminal where you start the API.

## Try it

With the server running, open another terminal.

Check the HTTP server:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

This endpoint checks that the HTTP server responds; it does not check the database on each request.

Upload a sample file:

```bash
echo "This is my first document." > /tmp/sample.txt
curl -i -F "file=@/tmp/sample.txt" http://localhost:8080/documents
```

Expect `201 Created` and JSON containing the document ID, original filename, `pending` status, and creation time. The internal storage path is left out of the response.

Check the saved file and database record:

```bash
ls -l uploads/
docker compose exec postgres psql -U document -d document_ai \
  -c "SELECT id, filename, status, storage_path FROM documents;"
```

## Project layout

```text
cmd/api/              Starts the server and connects the pieces
internal/handler/     Handles HTTP requests and responses
internal/service/     Handles file storage and processing logic
internal/repository/  Reads and writes database records
internal/model/       Defines the Go data structures
internal/database/    Sets up the database connection
migrations/           SQL files that create database tables
uploads/              Uploaded files, excluded from Git
```

## Common startup problems

- `DATABASE_URL is required`: set the environment variable in the terminal running Go.
- `address already in use`: another process is using the port. Stop your old server before restarting.
- `404` after adding a route: restart the Go server; `go run` does not reload code automatically.
- Database errors: check `docker compose logs --tail=50 postgres` and confirm your credentials and migrations.

## What comes next

1. Read uploaded text files and connect chunking to the upload flow.
2. Store the chunks and track processing status.
3. Generate embeddings, which represent text as numbers for finding related passages.
4. Retrieve relevant passages when someone asks a question.
5. Use an AI model to answer based on those passages.
