package chunks

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/tmc/langchaingo/textsplitter"
)

const (
	defaultSize    = 500
	defaultOverlap = 50
)

func SplitText(text string, size, overlap int) ([]string, error) {
	if size <= 0 {
		size = defaultSize
	}

	if overlap <= 0 || overlap >= size {
		overlap = defaultOverlap
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("empty text")
	}

	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(size),
		textsplitter.WithChunkOverlap(overlap),
		textsplitter.WithLenFunc(func(s string) int {
			return utf8.RuneCount([]byte(s)) //корректная обработки для русских символов, эмодзи и др.
		}),
	)

	chunks, err := splitter.SplitText(text)
	if err != nil {
		return nil, fmt.Errorf("error occured while splitting text: %w", err)
	}

	return chunks, nil

}
