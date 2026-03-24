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
4. Use `openplan list` and `openplan show <FULL_ID>` to inspect plan metadata and paths without dumping full bodies.
5. Use `openplan validate` to check the repository for invalid plans.

## Commands

- `openplan init [--prefix OPN]`
- `openplan new [title] [--tag TAG] [--parent ID]`
- `openplan edit <FULL_ID>`
- `openplan list [--all] [--status STATUS] [--tag TAG] [--json]`
- `openplan show <FULL_ID> [--json]`
- `openplan validate`
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
parent: ""
---

## Objective

## Context

## Research

## Plan

## Outputs

## Verification

## Review

## Notes
```

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
