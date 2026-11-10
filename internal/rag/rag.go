package rag

import (
	"context"
	"fmt"
	llms "ragExample/internal/llmthings"
	"ragExample/internal/repo"
	"strings"
)

type Service struct {
	repo   *repo.Repository
	client *llms.ClientLLM
	topK   int //топ чанков брать в контекст
}

type QueryRequest struct {
	Question string `json:"question"`
}

type QueryResult struct {
	Answer  string       `json:"answer"`
	Sources []repo.Chunk `json:"sources"`
}

func NewService(repository *repo.Repository, client *llms.ClientLLM, topK int) *Service {
	return &Service{
		repo:   repository,
		client: client,
		topK:   topK,
	}
}

// Данная функция запускает полный пайплайн RAG: вопрос -> эмбеддинги -> поиск -> генерация ответа
func (s *Service) Query(ctx context.Context, question string) (*QueryResult, error) {
	question = strings.TrimSpace(question)

	vectorQ, err := s.client.Embedder.EmbedQuery(ctx, question)
	if err != nil {
		return nil, fmt.Errorf("some error while embedding question: %w", err)
	}

	sources, err := s.repo.SearchSimilar(ctx, vectorQ, s.topK)
	if err != nil {
		return nil, fmt.Errorf("RAG Search: %w", err)
	}

	//дальше формируем промпт: вопрос полльзователя + системный промпт + контекст из sources

	allSources := buildContext(sources)
	finalPrompt := fmt.Sprintf(`
		Ты — ассистент по корпоративной базе знаний.
		Ответь на вопрос, используя ТОЛЬКО факты из контекста ниже.

		Правила:
		- Отвечай на русском языке.
		- Если вопрос короткий (название, термин) — собери из контекста всё, что относится к нему.
		- Если в контексте есть прямой ответ — дай его кратко и по существу.
		- Не выдумывай факты, которых нет в контексте.
		- Если информации недостаточно — ответь: «В документах недостаточно информации по этому вопросу.»

		Контекст:
		%s

		Вопрос: %s

		Ответ:
	`, allSources, question)

	//генерация
	answer, err := s.client.Generate(ctx, finalPrompt)
	if err != nil {
		return nil, fmt.Errorf("Generate answer: %w", err)
	}

	return &QueryResult{
		Answer:  answer,
		Sources: sources,
	}, nil

}

func buildContext(sources []repo.Chunk) string {
	var builder strings.Builder
	for i, v := range sources {
		builder.WriteString(fmt.Sprintf("------------\nДокумент №%d\n%s\n", i, v.Content))
	}

	return builder.String()
}
