package embeddings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type ClientLLM struct {
	llm      *ollama.LLM
	Embedder *embeddings.EmbedderImpl
}

func NewClientForEmbed(baseURL, modelForGeneration, modelForEmbeddings string) (*ClientLLM, error) {

	llm, err := ollama.New(ollama.WithModel(modelForGeneration), ollama.WithHTTPClient(&http.Client{}), ollama.WithServerURL(baseURL))

	if err != nil {
		return nil, fmt.Errorf("error while creating llm: %w", err)
	}

	llmForEmbed, err := ollama.New(ollama.WithModel(modelForEmbeddings))
	if err != nil {
		return nil, fmt.Errorf("error while creating llm: %w", err)
	}
	embedder, err := embeddings.NewEmbedder(llmForEmbed)

	if err != nil {
		return nil, fmt.Errorf("error while creating embedder: %w", err)
	}

	return &ClientLLM{
		llm:      llm,
		Embedder: embedder,
	}, nil
}

type embedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

func (c *ClientLLM) Embed(ctx context.Context, texts []string) (*embedResponse, error) {
	vectors, err := c.Embedder.EmbedDocuments(ctx, texts)

	if err != nil {
		return nil, fmt.Errorf("error while creating vectors: %w", err)
	}

	return &embedResponse{
		Embeddings: vectors,
	}, nil
}

func (c *ClientLLM) Generate(ctx context.Context, prompt string) (string, error) {
	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}
	resp, err := c.llm.GenerateContent(ctx, messages)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Content, nil
	}

	return "", fmt.Errorf("Модель не вернула ничего")
}
