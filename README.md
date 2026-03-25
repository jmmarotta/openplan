# OpenPlan

OpenPlan is a small Go CLI for filesystem-native technical planning. Each plan lives as a single markdown file in `.plans/`, stays reviewable in git, and remains easy to edit in a normal editor.

## Installation

Install the latest released version with Go:

```sh
go install github.com/jmmarotta/openplan/cmd/openplan@latest
```

Install the local checkout from the repository root:

```sh
go install ./cmd/openplan
```

## Goals

- Keep plans local and on disk.
- Make plan changes easy to review in git.
- Keep the CLI thin and deterministic.
- Leave room for future editor integration without adding premature abstraction.

## Workflow

1. Run `openplan init --prefix OPN` in a repository root.
2. Run `openplan new "Draft README"` to create a templated plan and open it in `$EDITOR`.
3. Edit frontmatter and body sections directly in the markdown file under `.plans/`.
4. Use `openplan list` or `openplan ls`, and `openplan show <FULL_ID>`, to inspect plan metadata and paths without dumping full bodies.
5. Use `openplan validate` to check the repository for invalid plans.

## Commands

- `openplan init` (`openplan i`) `[--prefix OPN]`
- `openplan new` (`openplan n`) `[title] [--tag TAG] [--parent ID]`
- `openplan edit` (`openplan e`) `<FULL_ID>`
- `openplan list` (`openplan ls`) `[--all] [--status STATUS] [--tag TAG] [--json]`
- `openplan show <FULL_ID> [--json]`
- `openplan validate` (`openplan v`)
- `openplan skill`

`openplan init` creates `.plans/openplan.jsonc` by default. `openplan.json` is also accepted.

## Plan Format

Plan files live in `.plans/` and use filenames like `.plans/OPN-12_7K4M9XQ2.md`.

```md
---
id: OPN-12_7K4M9XQ2
title: "Draft README"
status: inbox # inbox, plan, active, done, cancelled
tags:
  - docs
---

## Objective

## Context

## Research

## Open Questions

## System Surfaces

## Invariants

## Outputs

## Verification

## Execution Plan

## Notes
```

`Open Questions` should always be present. Use `None.` when there are no open
questions to record.

`System Surfaces` should describe the main caller-facing interfaces and
entrypoints in an interface-first way, grouped by file or ownership boundary.

Section authoring conventions:

- Keep the canonical top-level sections flat. Add nested subsections only when a section grows large enough to need internal structure.
- In `System Surfaces`, use a file or ownership-boundary header first, then list the meaningful surfaces under it.
- Default each surface entry to a code block containing an APOSD-style interface comment plus the signature, command, or surface line itself.
- Use extra bullets under a surface only when the code block alone cannot capture an important boundary, constraint, or dependency.
- If a low-level file has no meaningful public surface, document the higher-level surface instead of inventing a shallow one.
- Use `Invariants` for the non-negotiable ownership rules, interface guarantees, and design constraints that must remain true through implementation.

Example `System Surfaces` entry:

````md
## System Surfaces

### `lua/openplan/cli.lua`

#### `list_plans(opts?)`

```lua
--- Return the current repository's OpenPlan catalog in a plugin-ready shape.
--- This surface centralizes CLI invocation and decoding so callers do not depend
--- on process execution details or wire-format knowledge.
---@param opts? { all?: boolean }
---@return { plans: openplan.PlanRow[], issues: openplan.ValidationIssue[] }|nil, string? err
list_plans(opts)
```
````

## ID Semantics

- IDs use `<PREFIX>-<NUMBER>_<SUFFIX>`.
- The full ID is the authoritative identifier in v1.
- `<NUMBER>` is allocated from the local repository state.
- `<SUFFIX>` is an immutable 8-character Crockford base32 string.

## Verification

- `go test ./...`
- `golangci-lint run`

## Notes

- OpenPlan does not require git for correctness, but git improves reviewability.
- `list` keeps valid plans visible even when invalid files exist, and reports issues separately.
- `show` reports validation failures for the target plan instead of printing a misleading partial summary.
- `tags` and `parent` are optional in frontmatter and can be omitted when empty.
