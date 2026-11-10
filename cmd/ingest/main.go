package ingest

import (
	"context"
	"fmt"
	"os"
	"ragExample/internal/chunks"
	"ragExample/internal/config"
	"ragExample/internal/documentloader"
	"ragExample/internal/repo"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()
	entries, err := os.ReadDir("../internal/docs")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rep := repo.NewRepository(cfg.DatabaseURL)

	llmForEmbed, err := ollama.New(ollama.WithModel("bge-m3"))
	if err != nil {
		fmt.Printf("error while creating llm: %w", err)
		return
	}
	embedder, err := embeddings.NewEmbedder(llmForEmbed)

	if err != nil {
		fmt.Printf("error while creating embedder: %w", err)
		return
	}
	for _, v := range entries {
		DocID, err := rep.AddDocument(ctx, v.Name(), "/docs"+v.Name())
		if err != nil {
			fmt.Printf("ERROR add doc: %s", err.Error())
			continue
		}
		DocumentString, err := documentloader.LoadDocument(v.Name(), ctx)
		if err != nil {
			fmt.Printf("skip file %s\n", v.Name())
			continue
		}

		chunks, err := chunks.SplitText(DocumentString, 0, 0)
		if err != nil {
			fmt.Printf("error in chunking: %s\n", v.Name())
			continue
		}

		embedds, err := embedder.EmbedDocuments(ctx, chunks)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		rep.AddChunk(ctx)
	}

}
