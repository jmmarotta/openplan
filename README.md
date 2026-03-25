# OpenPlan

OpenPlan is a CLI for managing technical plans as markdown files in your codebase.

Plans stay next to the code they describe, easy to open in your editor and work
through with your agent. The [openplan.nvim](https://github.com/jmmarotta/openplan.nvim)
plugin makes this seamless from Neovim.

## Why OpenPlan

Planning tools usually live outside your editor -- in a browser tab, a separate
app, or a thread you have to go find later. Planning directly in an agent session
is convenient but hard to review, not persisted, and makes it difficult to
reference, edit, or track the status of each plan.

OpenPlan keeps plans where the work happens: in your repository, in your editor,
alongside the code they describe.

- Plans are markdown files in `.plans/`, so they're easy to read, edit, and review.
- The CLI creates, validates, and lists plans without hiding them behind a database.
- Editor plugins like [openplan.nvim](https://github.com/jmmarotta/openplan.nvim) and [opencode.nvim](https://github.com/nickjvandyke/opencode.nvim) let you plan with your agent without leaving your workflow.

## Installation

Install the latest released version with Go:

```sh
go install github.com/jmmarotta/openplan/cmd/openplan@latest
```

Install from a local checkout:

```sh
go install ./cmd/openplan
```

## Quick Start

Initialize OpenPlan in a repository:

```sh
openplan init --prefix OPN
```

Create a new plan:

```sh
openplan new "Draft README"
```

List plans:

```sh
openplan list
```

Validate the plan directory:

```sh
openplan validate
```

`openplan init` creates `.plans/openplan.jsonc` by default. `openplan.json` is
also accepted.

## Commands

- `openplan init` (`openplan i`) `[--prefix OPN]`: create OpenPlan config for the repo
- `openplan new` (`openplan n`) `[title] [--tag TAG] [--parent ID] [--no-edit] [--json]`: create a templated plan, optionally skip opening `$EDITOR`, and optionally return machine-readable metadata
- `openplan edit` (`openplan e`) `<FULL_ID>`: open an existing plan in `$EDITOR`
- `openplan list` (`openplan ls`) `[--all] [--status STATUS] [--tag TAG] [--json]`: list plans and validation issues
- `openplan show` `<FULL_ID> [--json]`: show metadata for one plan without printing the full body
- `openplan validate` (`openplan v`): validate all plans in the repository
- `openplan skill`: print the bundled OpenPlan skill text

## Plan Files

Plans live in `.plans/` and use filenames like `.plans/OPN-12_7K4M9XQ2.md`.
Each file is standard markdown with YAML frontmatter.

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

## Touch Points

## Outputs

## Verification

## Execution Plan

## Notes
```

`tags` and `parent` are optional and can be omitted when empty.

## Template Sections

OpenPlan's default template is meant to keep plans concrete and reviewable:

- `Objective`: what success looks like
- `Context`: current state, constraints, and background
- `Research`: facts, references, or exploration notes
- `Open Questions`: unresolved decisions or unknowns; use `None.` when there are none
- `System Surfaces`: the main interfaces, commands, APIs, or entrypoints affected by the work
- `Invariants`: rules or guarantees that must remain true
- `Touch Points`: the concrete files and artifacts expected to change; use `None.` when there are none
- `Outputs`: the expected artifacts or deliverables
- `Verification`: how you will confirm the work is correct
- `Execution Plan`: the ordered implementation approach
- `Notes`: anything useful that does not fit the main structure

Section authoring conventions:

- `Open Questions` should always be present. Use `None.` when there are no open questions to record.
- `Touch Points` should always be present. Use `None.` when there are no concrete artifacts to record.
- `System Surfaces` should describe the main caller-facing interfaces and entrypoints in an interface-first way.
- `Invariants` should capture the rules and guarantees that must remain true through implementation.

## IDs and Behavior

- IDs use the format `<PREFIX>-<NUMBER>_<SUFFIX>`.
- The full ID is the authoritative identifier.
- `<NUMBER>` is allocated from repository-local state.
- `<SUFFIX>` is an immutable 8-character Crockford base32 string.
- `list` keeps valid plans visible even when invalid files exist, and reports issues separately.
- `show` reports validation failures for the target plan instead of printing a partial or misleading summary.

OpenPlan does not require git for correctness, but git makes plan review much
more useful in practice.

## Development

For local verification:

```sh
go test ./...
golangci-lint run
```

Install `lefthook` locally so the repository hooks can run.
