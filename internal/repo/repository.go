package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func formatEmbed(embedding []float32) string {
	parts := make([]string, len(embedding))
	for i := 0; i < len(parts); i++ {
		parts[i] = fmt.Sprintf("%g", embedding[i])
	}

	result := "[" + strings.Join(parts, ", ") + "]"
	return result
}

func NewRepository(connString string) *Repository {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil
	}

	return &Repository{
		pool: pool,
	}
}

type Document struct {
	ID         int64
	Title      string
	SourcePath string
	CreatedAt  time.Time
}

type Chunk struct {
	ID         int64
	DocID      int64
	Content    string
	ChunkIndex int
	Score      float64
}

func (r *Repository) AddDocument(ctx context.Context, title, sourcePath string) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx, `INSERT INTO documents (title, source_path) VALUES $1, $2 returning id`, title, sourcePath).Scan(&id)

	return id, err
}

func (r *Repository) AddChunk(ctx context.Context, content string, docID int64, chunkIndex int, embedding []float32) error {
	vec := formatEmbed(embedding)
	_, err := r.pool.Exec(ctx, `insert into chunks (document_id, content, chunk_index, embedding) values $1, $2, $3, $4::vector`, docID, content, chunkIndex, vec)
	return err
}

func (r *Repository) ListDocuments(ctx context.Context) ([]Document, error) {
	rows, err := r.pool.Query(ctx, `select id, title, source_path, created_at from documents`)

	if err != nil {
		return nil, fmt.Errorf("listDocs: error while fetching rows: %w", err)
	}

	defer rows.Close()

	var docs []Document
	for rows.Next() {
		var doc Document
		err := rows.Scan(&doc.ID, &doc.Title, &doc.SourcePath, &doc.CreatedAt)
		if err != nil {
			fmt.Println("error while fetching row, skipping")
		}

		docs = append(docs, doc)
	}
	return docs, nil
}

func (r *Repository) SearchSimilar(ctx context.Context, embedding []float32, limit int) ([]Chunk, error) {
	vec := formatEmbed(embedding)
	rows, err := r.pool.Query(ctx, `
		select c.id, c.document_id, c.content, c.chunk_index, 1 - (c.embedding <=> $1::vector) as score
		from chunks c
		where c.embedding is not null
		order by c.embedding <=> $1::vector desc
		limit $2
	`, vec, limit)

	if err != nil {
		return nil, fmt.Errorf("SearchSimilar: error while fetching rows: %w", err)
	}
	defer rows.Close()

	var chunks []Chunk
	for rows.Next() {
		var chunk Chunk

		err := rows.Scan(&chunk.ID, &chunk.DocID, &chunk.Content, &chunk.ChunkIndex, &chunk.Score)
		if err != nil {
			fmt.Println("chunks read: error skipping 1 row")
		}

		chunks = append(chunks, chunk)
	}

	return chunks, nil
}
