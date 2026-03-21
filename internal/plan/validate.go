package plan

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var tagPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

func Validate(doc Document) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	addIssue := func(field string, message string) {
		issues = append(issues, ValidationIssue{Path: doc.Path, Field: field, Message: message})
	}

	// Validate ID
	if doc.Frontmatter.ID == "" {
		addIssue("id", "is required")
	} else {
		if _, err := ParseID(doc.Frontmatter.ID); err != nil {
			addIssue("id", "must be a valid full plan ID")
		}

		stem := strings.TrimSuffix(filepath.Base(doc.Path), filepath.Ext(doc.Path))
		if stem != "" && stem != doc.Frontmatter.ID {
			addIssue("id", "must match the filename stem")
		}
	}

	// Validate title
	if strings.TrimSpace(doc.Frontmatter.Title) == "" {
		addIssue("title", "is required")
	}

	// Validate status
	if !doc.Frontmatter.Status.Valid() {
		addIssue("status", "must be one of inbox, plan, active, done, or cancelled")
	}

	// Validate tags
	if doc.Frontmatter.Tags == nil {
		addIssue("tags", "is required")
	} else {
		seenTags := make(map[string]struct{}, len(doc.Frontmatter.Tags))
		for _, tag := range doc.Frontmatter.Tags {
			normalized := normalizeTag(tag)
			// Keep tag reporting field-oriented and deterministic instead of
			// emitting one issue per element.
			if tag != normalized || !tagPattern.MatchString(normalized) {
				addIssue("tags", "must contain normalized lowercase tags")
				break
			}
			if _, ok := seenTags[tag]; ok {
				addIssue("tags", "must not contain duplicates")
				break
			}
			seenTags[tag] = struct{}{}
		}
	}

	// Validate parent
	if doc.Frontmatter.Parent != "" {
		if _, err := ParseID(doc.Frontmatter.Parent); err != nil {
			addIssue("parent", "must be empty or a valid full plan ID")
		}
	}

	return issues
}

func NormalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return []string{}
	}

	normalized := make([]string, 0, len(tags))
	seen := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		value := normalizeTag(tag)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}
	sort.Strings(normalized)
	return normalized
}

func normalizeTag(tag string) string {
	return strings.ToLower(strings.TrimSpace(tag))
}
