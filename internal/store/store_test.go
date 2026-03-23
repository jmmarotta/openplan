package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"openplan/internal/plan"
	"openplan/internal/repo"
)

func TestCreateGetAndList(t *testing.T) {
	root := t.TempDir()
	ctx, err := repo.Init(root, repo.DefaultConfig("OPN"))
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	seedPlan(t, ctx, plan.Frontmatter{ID: "OPN-2_ZZZZZZZZ", Title: "Done plan", Status: plan.StatusDone, Tags: []string{"docs"}})
	seedPlan(t, ctx, plan.Frontmatter{ID: "OPN-2_AAAAAAAA", Title: "Active plan", Status: plan.StatusActive, Tags: []string{"cli"}})
	seedPlan(t, ctx, plan.Frontmatter{ID: "OPN-10_BBBBBBBB", Title: "Planning work", Status: plan.StatusPlan, Tags: []string{"backend"}})
	if err := os.WriteFile(filepath.Join(ctx.PlansDir, "broken.md"), []byte("oops"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	created, err := Create(ctx, NewPlanInput{Title: "New task", Tags: []string{"CLI", "feature"}})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	parsedID, err := plan.ParseID(created.Frontmatter.ID)
	if err != nil {
		t.Fatalf("created ID parse returned error: %v", err)
	}
	if parsedID.Number != 11 {
		t.Fatalf("created ID number = %d, want %d", parsedID.Number, 11)
	}

	got, err := Get(ctx, created.Frontmatter.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Document == nil {
		t.Fatalf("Get returned nil document")
	}
	if got.Document.Frontmatter.Title != "New task" {
		t.Fatalf("Get title = %q, want %q", got.Document.Frontmatter.Title, "New task")
	}

	result, err := List(ctx, Query{})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(result.Issues) != 1 {
		t.Fatalf("List issues = %d, want 1", len(result.Issues))
	}
	if len(result.Documents) != 3 {
		t.Fatalf("List documents = %d, want 3", len(result.Documents))
	}
	if result.Documents[0].Frontmatter.ID != "OPN-2_AAAAAAAA" {
		t.Fatalf("first listed plan = %q, want %q", result.Documents[0].Frontmatter.ID, "OPN-2_AAAAAAAA")
	}
}

func TestCreateRejectsInvalidParent(t *testing.T) {
	root := t.TempDir()
	ctx, err := repo.Init(root, repo.DefaultConfig("OPN"))
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	_, err = Create(ctx, NewPlanInput{Title: "Child", Parent: "OPN-9_UNKNOWN1"})
	if !errors.Is(err, errInvalidParent) {
		t.Fatalf("Create error = %v, want wrapped invalid parent error", err)
	}
}

func TestGetInvalidPlan(t *testing.T) {
	root := t.TempDir()
	ctx, err := repo.Init(root, repo.DefaultConfig("OPN"))
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	path := ctx.PlanPath("OPN-1_ABCDEFGH")
	if err := os.WriteFile(path, []byte("---\nid: wrong\n---\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	result, err := Get(ctx, "OPN-1_ABCDEFGH")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if len(result.ValidationIssues) == 0 {
		t.Fatalf("Get validation issues = %d, want > 0", len(result.ValidationIssues))
	}
}

func seedPlan(t *testing.T, ctx repo.Context, meta plan.Frontmatter) {
	t.Helper()
	meta.Tags = plan.NormalizeTags(meta.Tags)
	content := plan.Template(meta)
	if err := os.WriteFile(ctx.PlanPath(meta.ID), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
