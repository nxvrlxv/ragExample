package documentloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tmc/langchaingo/documentloaders"
)

var (
	folder = "../docs"
)

func LoadDocument(filename string, ctx context.Context) (string, error) {
	fullPath := filepath.Join(folder, filename)
	f, err := os.Open(fullPath)
	if err != nil {
		return "", fmt.Errorf("Error to open file")
	}

	loader := documentloaders.NewText(f)
	document, err := loader.Load(ctx)
	if err != nil {
		return "", fmt.Errorf("Error load: %w", err)
	}

	return document[0].PageContent, nil

	//TODO: сделать формирование строки из списка Document через strings.Builder
}
