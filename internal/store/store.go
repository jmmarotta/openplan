package store

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/jmmarotta/openplan/internal/plan"
	"github.com/jmmarotta/openplan/internal/repo"
)

// Query describes the filters applied when listing plans from disk
type Query struct {
	IncludeClosed bool
	Statuses      []plan.Status
	Tag           string
}

// ListResult contains valid plan documents plus any issues found while loading
// other plan files
type ListResult struct {
	Documents []plan.Document        `json:"documents"`
	Issues    []plan.ValidationIssue `json:"issues,omitempty"`
}

// NewPlanInput describes the user-editable metadata needed to create a plan
type NewPlanInput struct {
	Title  string
	Tags   []string
	Parent string
}

// GetResult contains the outcome of loading a single plan path. Validation
// issues are reported as data so callers can decide how to surface them.
type GetResult struct {
	Path             string                 `json:"path"`
	Document         *plan.Document         `json:"document,omitempty"`
	ValidationIssues []plan.ValidationIssue `json:"validationIssues,omitempty"`
}

var (
	errPlanNotFound  = errors.New("plan not found")
	errInvalidParent = errors.New("invalid parent plan")
)

// List loads plan files, validates them, applies filters, and returns valid
// documents alongside any file issues that were encountered
func List(ctx repo.Context, q Query) (ListResult, error) {
	paths, err := planFiles(ctx)
	if err != nil {
		return ListResult{}, err
	}

	result := ListResult{}
	for _, path := range paths {
		doc, err := plan.ParseFile(path)
		if err != nil {
			result.Issues = append(result.Issues, plan.ValidationIssue{Path: path, Field: "frontmatter", Message: err.Error()})
			continue
		}

		issues := plan.Validate(doc)
		if len(issues) > 0 {
			result.Issues = append(result.Issues, issues...)
			continue
		}

		if !matchesQuery(doc, q) {
			continue
		}
		result.Documents = append(result.Documents, doc)
	}

	sortDocuments(result.Documents)
	sortIssues(result.Issues)
	return result, nil
}

// Get loads a single plan by full ID. Validation failures are returned in the
// result so callers can surface repairable files without depending on an error
// type.
func Get(ctx repo.Context, fullID string) (GetResult, error) {
	if _, err := plan.ParseID(fullID); err != nil {
		return GetResult{}, err
	}

	path := ctx.PlanPath(fullID)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return GetResult{}, fmt.Errorf("%w %q", errPlanNotFound, fullID)
		}
		return GetResult{}, err
	}

	doc, err := plan.ParseBytes(path, data)
	if err != nil {
		return GetResult{
			Path: path,
			ValidationIssues: []plan.ValidationIssue{{
				Path:    path,
				Field:   "frontmatter",
				Message: err.Error(),
			}},
		}, nil
	}

	issues := plan.Validate(doc)
	if len(issues) > 0 {
		return GetResult{Path: path, ValidationIssues: issues}, nil
	}

	return GetResult{Path: path, Document: &doc}, nil
}

// Create allocates the next plan ID, writes the templated file, and returns
// the parsed document for the content that was created
func Create(ctx repo.Context, input NewPlanInput) (doc plan.Document, err error) {
	parent := strings.TrimSpace(input.Parent)
	if parent != "" {
		if _, err := plan.ParseID(parent); err != nil {
			return plan.Document{}, fmt.Errorf("%w %q: must be a valid full plan ID", errInvalidParent, parent)
		}
		parentResult, err := Get(ctx, parent)
		if err != nil {
			return plan.Document{}, fmt.Errorf("%w %q: %w", errInvalidParent, parent, err)
		}
		if len(parentResult.ValidationIssues) > 0 || parentResult.Document == nil {
			return plan.Document{}, fmt.Errorf("%w %q: must refer to a valid plan", errInvalidParent, parent)
		}
	}

	existingIDs, err := existingIDs(ctx)
	if err != nil {
		return plan.Document{}, err
	}

	newID, err := plan.NewID(ctx.Config.Prefix, existingIDs, rand.Reader)
	if err != nil {
		return plan.Document{}, err
	}

	meta := plan.Frontmatter{
		ID:     plan.FormatID(newID),
		Title:  defaultTitle(input.Title),
		Status: plan.StatusInbox,
		Tags:   plan.NormalizeTags(input.Tags),
		Parent: parent,
	}
	path := ctx.PlanPath(meta.ID)

	preview := plan.Document{Path: path, Frontmatter: meta}
	issues := plan.Validate(preview)
	if len(issues) > 0 {
		return plan.Document{}, fmt.Errorf("invalid plan metadata: %s: %s", issues[0].Field, issues[0].Message)
	}

	content := plan.Template(meta)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return plan.Document{}, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if _, err := file.WriteString(content); err != nil {
		return plan.Document{}, err
	}

	doc, err = plan.ParseBytes(path, []byte(content))
	return doc, err
}

// defaultTitle keeps plan creation editor-first by supplying a stable fallback
// when the user omits a title
func defaultTitle(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return "Untitled plan"
	}
	return title
}

// existingIDs collects valid filename-based plan IDs so Create can allocate the
// next repo-local number without trusting invalid files
func existingIDs(ctx repo.Context) ([]string, error) {
	paths, err := planFiles(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(paths))
	for _, path := range paths {
		stem := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		if _, err := plan.ParseID(stem); err == nil {
			ids = append(ids, stem)
		}
	}
	return ids, nil
}

// planFiles lists markdown plan files directly under `.plans/` in deterministic
// order so higher-level operations can layer validation and filtering on top
func planFiles(ctx repo.Context) ([]string, error) {
	entries, err := os.ReadDir(ctx.PlansDir)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		paths = append(paths, filepath.Join(ctx.PlansDir, entry.Name()))
	}

	sort.Strings(paths)
	return paths, nil
}

// matchesQuery applies the list filters to a single validated document
func matchesQuery(doc plan.Document, q Query) bool {
	if len(q.Statuses) > 0 {
		if !slices.Contains(q.Statuses, doc.Frontmatter.Status) {
			return false
		}
	} else if !q.IncludeClosed && (doc.Frontmatter.Status == plan.StatusDone || doc.Frontmatter.Status == plan.StatusCancelled) {
		return false
	}

	if q.Tag == "" {
		return true
	}

	return slices.Contains(doc.Frontmatter.Tags, q.Tag)
}

// sortDocuments keeps list output deterministic by ordering plans by numeric ID
// first and suffix second, with the full ID as a final tie-breaker
func sortDocuments(docs []plan.Document) {
	sort.Slice(docs, func(i int, j int) bool {
		left, leftErr := plan.ParseID(docs[i].Frontmatter.ID)
		right, rightErr := plan.ParseID(docs[j].Frontmatter.ID)
		if leftErr != nil || rightErr != nil {
			return docs[i].Frontmatter.ID < docs[j].Frontmatter.ID
		}
		if left.Number != right.Number {
			return left.Number < right.Number
		}
		if left.Suffix != right.Suffix {
			return left.Suffix < right.Suffix
		}
		return docs[i].Frontmatter.ID < docs[j].Frontmatter.ID
	})
}

// sortIssues keeps issue reporting stable across runs and filesystem ordering
func sortIssues(issues []plan.ValidationIssue) {
	sort.Slice(issues, func(i int, j int) bool {
		left := strings.Join([]string{issues[i].Path, issues[i].Field, issues[i].Message}, "\x00")
		right := strings.Join([]string{issues[j].Path, issues[j].Field, issues[j].Message}, "\x00")
		return left < right
	})
}
