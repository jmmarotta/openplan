Create a new repository called `OpenPlan`.

Build a local-first planning tool inspired by `seal`, but much lighter: one markdown file per plan, stored directly in the repo filesystem, with a Go CLI built using Cobra.

Purpose:
OpenPlan is a filesystem-native planning pad and issue tracker for technical work. It should support a research -> plan -> implement workflow inside a single markdown file per plan, make plans easy to review in git and Neovim, and keep all planning artifacts on disk instead of rendering plans back to the user in chat.

Core design requirements:
- Use a hidden repo-local directory: `.plans/`
- Store each plan as a single markdown file with YAML frontmatter
- Use canonical plan IDs with this format: `<PREFIX>-<NUMBER>_<UNIQUE_SUFFIX>`
- Use filenames like `.plans/OPN-12_7K4M9XQ2.md`
- `openplan init --prefix OPN` configures the ticket prefix
- Treat this as repo-user/local-first for v1: multiple repo users may create plans in separate clones and merge later
- Do not maintain a separate index file; generate views live from frontmatter
- Do not output full plan bodies to the user; write and edit plans in the filesystem only

ID rules:
- The full ID is the authoritative plan identifier
- `<NUMBER>` is a short, human-friendly local sequence allocated from the current repo state at creation time
- `<NUMBER>` is not guaranteed to remain globally unique across merged histories and may repeat after merges
- `<UNIQUE_SUFFIX>` exists to make IDs merge-friendlier and uniquely distinguish plans created with the same `<NUMBER>` in different clones or branches
- Encode `<UNIQUE_SUFFIX>` as a permanent 8-character random Crockford base32 string
- The suffix is immutable once assigned
- The suffix exists for uniqueness, not for human-readable timestamps
- Commands may accept `<NUMBER>` as shorthand only when it resolves to exactly one plan; otherwise they must require the full ID

Plan lifecycle:
- `inbox`
- `plan`
- `active`
- `done`
- `cancelled`

Plan file format:
Each plan should be a single markdown document with YAML frontmatter and a fully templated body.

Required frontmatter:
- `id`
- `title`
- `status`
- `tags`
- `parent`

Body template:
## Objective

## Context

## Research

## Open Questions

## System Surfaces

## Invariants

## Touch Points

## Outputs

## Verification

## Execution Plan

## Notes

Planning conventions:
- `Open Questions` should always exist; use `None.` when there are no unresolved questions.
- `Touch Points` should always exist; use `None.` when there are no concrete artifacts to record.
- Keep the canonical top-level sections flat. Add nested subsections only when a section grows large enough to need internal structure.
- `System Surfaces` should be interface-first and grouped by file or ownership boundary.
- Under each file or ownership boundary, list the meaningful surfaces rather than every helper.
- If a file has multiple meaningful surfaces, prefer light bullet labels over extra headers before the code blocks.
- Prefer a code block with an APOSD-style interface comment plus the signature, command, or surface line for each `System Surfaces` entry.
- Add bullets only when the code block alone cannot explain an important boundary, constraint, or dependency.
- If a low-level file has no meaningful public surface, document the higher-level surface instead of inventing a shallow one.
- Use `Invariants` for non-negotiable ownership rules, interface guarantees, and design constraints that must remain true during implementation.
- Use `Touch Points` as a compact artifact inventory, not a second design section.
- Group `Touch Points` by repo or subsystem when useful.
- Format each touch point as a bullet with `A`, `M`, `D`, or `R`, with the path bolded on the first line and one short explanation below it.

Behavior:
- `new` creates a fully templated file and immediately opens it in `$EDITOR`
- `edit <ID>` opens an existing plan in `$EDITOR`
- plan status is edited directly in frontmatter rather than through a dedicated `set-status` command
- plan title lives in frontmatter only; do not duplicate it as a body H1
- edits should not rewrite, reorder, or move body sections unless explicitly requested
- `parent` is used for followups, derivation, and duplicate lineage instead of introducing a separate `duplicate` status
- cancelled/duplicate rationale can live in `Notes`

Validation invariants:
- filename stem must exactly match frontmatter `id`
- `status` must be one of the defined lifecycle values
- `tags` must be a list of normalized lowercase strings
- `parent` must be empty or a valid plan ID
- files with invalid frontmatter should be surfaced clearly by read commands

CLI commands for v1:
- `init`: initialize `.plans/` and repo config, with optional ticket prefix
- `new`: create the next plan file from template and open it in editor
- `edit`: open a plan by ID
- `list`: generate a table view from frontmatter
- `show`: show plan metadata, path, or a concise summary without dumping the full plan body
- `skill`: print deterministic Markdown guidance for AI agents describing how to use OpenPlan correctly

Lookup behavior:
- `edit` and `show` should accept the full canonical ID
- `edit` and `show` may accept the numeric portion as shorthand only when it resolves uniquely
- ambiguous numeric lookups should fail with a clear message that lists the matching full IDs
- merged histories may contain multiple plans with the same numeric portion, so the full ID remains the only stable identifier
- if users want numeric portions to be unique again after a merge, they must reconcile those numbers manually outside v1 tooling

List behavior:
- default `list` should exclude `done` and `cancelled`
- support `--all`
- support filtering by `--status`
- support filtering by `--tag`
- support `--json` for machine-readable output
- sort list results by numeric portion first, with a deterministic secondary tie-break using the suffix
- output should be friendly for terminal and future Neovim integration

Show behavior:
- do not dump the full markdown body by default
- support `--json`
- include enough metadata for editor integrations to resolve and open the file cleanly

Skill behavior:
- `skill` should print a deterministic Markdown document to stdout
- it should explain the OpenPlan workflow, commands, statuses, file format, ID semantics, and how AI agents should interact with plan files
- model the UX after the `playwriter` skill command, but scoped to OpenPlan

Implementation expectations:
- Use Go
- Use Cobra for the CLI
- Keep the architecture small, obvious, and easy to extend
- Preserve strong filesystem ergonomics and git reviewability
- Prefer simple, explicit data modeling over abstraction
- Keep the codebase ready for a future Neovim plugin, but do not implement the plugin yet
- Design the code so editor integration can be layered on later
- Keep command handlers thin and centralize plan parsing, validation, ID generation, and rendering logic in focused packages

Repository expectations:
- Include a clear README with goals, workflow, commands, file format, ID format, and examples
- Include sensible project structure for a small Go CLI
- Include tests for plan parsing, ID generation, template creation, shorthand ID resolution, status validation, and list filtering
- Include fixtures for valid and invalid plan files if useful
- Include linting with a sensible Go setup such as `golangci-lint`
- Use `lefthook`
- Test `skill` output with a golden file or equivalent deterministic snapshot

Workflow expectations for the agent building this repo:
- Research existing patterns as needed, especially `seal`, Cobra CLI structure, markdown + YAML frontmatter handling in Go, and the `playwriter` skill command
- Produce an internal plan in the filesystem while working, not in chat
- Implement incrementally
- Verify the implementation with tests
- Keep chat output concise and focused on outcomes, files changed, and verification results

Non-goals for v1:
- no Neovim plugin yet
- no server
- no remote sync service
- no separate database
- no generated index file checked into git
- no complex workflow engine beyond the defined statuses
- no requirement that the numeric portion of an ID be globally unique across merged histories

Use `seal` as the main reference for product feel:
- `/Users/julianmarotta/projects/seal`
- filesystem-native workflow
- strong reviewability in git
- explicit command ergonomics
- structured documents
But simplify it into a single-file-per-plan model.

The repository should feel polished, minimal, practical, and merge-friendly.
