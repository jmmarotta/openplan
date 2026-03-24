package plan

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseFile(path string) (Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}
	return ParseBytes(path, data)
}

func ParseBytes(path string, data []byte) (Document, error) {
	const openingDelimiter = "---\n"
	const closingDelimiter = "\n---\n"

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(text, openingDelimiter) {
		return Document{}, fmt.Errorf("missing YAML frontmatter")
	}

	frontmatterText, body, ok := strings.Cut(text[len(openingDelimiter):], closingDelimiter)
	if !ok {
		return Document{}, fmt.Errorf("unterminated YAML frontmatter")
	}

	var fm Frontmatter
	decoder := yaml.NewDecoder(strings.NewReader(frontmatterText))
	decoder.KnownFields(true)
	if err := decoder.Decode(&fm); err != nil {
		return Document{}, fmt.Errorf("parse frontmatter: %w", err)
	}
	if fm.Tags == nil {
		fm.Tags = []string{}
	}

	return Document{
		Path:        path,
		Frontmatter: fm,
		Body:        strings.TrimPrefix(body, "\n"),
	}, nil
}
