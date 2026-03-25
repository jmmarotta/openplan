---
name: openplan
description: Use for OpenPlan repositories when you need to initialize planning, list plans, create new plans, or edit `.plans/` files directly.
---

# OpenPlan Skill

Use OpenPlan as a filesystem-native planning system that keeps the source of truth on disk.

## Workflow

- Use `openplan init` to initialize OpenPlan in a repository.
- Use `openplan list` to discover current work.
- Use `openplan new [title]` to create a new plan file.
- Use `openplan --help` when you need the full command surface.
- Read, edit, update, or delete plan files with filesystem tools instead of asking the CLI to mutate plan contents.

## Statuses

- `inbox`: captured but not yet shaped.
- `plan`: researched and scoped.
- `active`: implementation in progress.
- `done`: completed and verified.
- `cancelled`: stopped or superseded.

## File Format

- Plans live in `.plans/`.
- Each plan is a single markdown file with YAML frontmatter.
- Required frontmatter fields: `id`, `title`, `status`, `tags`, `parent`.
- The canonical body sections are `Objective`, `Context`, `Research`, `Open Questions`, `System Surfaces`, `Invariants`, `Touch Points`, `Outputs`, `Verification`, `Execution Plan`, and `Notes`.
- `Open Questions` should always be present. Use `None.` when there are no open questions to record.
- `Touch Points` should always be present. Use `None.` when there are no concrete artifacts to record.
- Keep the canonical top-level sections flat. Add nested subsections only when a section grows large enough to need internal structure.
- `System Surfaces` should be interface-first: group by file or ownership boundary, then list the meaningful surfaces under it.
- If a file has multiple meaningful surfaces, prefer light bullet labels over extra headers before the code blocks.
- Default each `System Surfaces` entry to a code block with an APOSD-style interface comment plus the signature, command, or surface line itself.
- Add bullets below a `System Surfaces` entry only when the code block alone is not enough to explain an important boundary, constraint, or dependency.
- If a low-level file has no meaningful public surface, document the higher-level surface instead of inventing a shallow one.
- Use `Invariants` for non-negotiable ownership rules, interface guarantees, and design constraints that must remain true during implementation.
- Use `Touch Points` as a compact artifact inventory, not a second design section.
- Group `Touch Points` by repo or subsystem when useful.
- Format each touch point as a bullet with `A`, `M`, `D`, or `R`, with the path bolded on the first line and one short explanation below it.

Example `System Surfaces` shape:

````md
## System Surfaces

### `lua/openplan/cli.lua`

- `list_plans(opts?)`

```lua
--- Return the current repository's OpenPlan catalog in a plugin-ready shape.
--- This surface centralizes CLI invocation and decoding so callers do not depend
--- on process execution details or wire-format knowledge.
---@param opts? { all?: boolean }
---@return { plans: openplan.PlanRow[], issues: openplan.ValidationIssue[] }|nil, string? err
list_plans(opts)
```
````

Example `Touch Points` shape:

```md
## Touch Points

### `openplan.nvim`

- **A `lua/openplan/cli.lua`**
  Add the machine-facing adapter for list and create flows.

- **M `README.md`**
  Document installation, commands, and picker behavior.
```

## ID Rules

- Full IDs are authoritative in v1.
- IDs use the format `<PREFIX>-<NUMBER>_<SUFFIX>`.
- `<NUMBER>` is allocated from the local repository state.
- `<SUFFIX>` is an immutable 8-character Crockford base32 string that keeps IDs merge-friendly.

## Agent Guidance

- Do not treat chat output as the plan body; write planning artifacts to plan files.
- Do not rewrite or reorder sections unless the user asks.
- Prefer filesystem reads over CLI rendering when you need plan contents.
- Use `list` to discover plans and `new` to allocate IDs, then work on files directly.
- Treat invalid plans as files to repair in the editor, not as data to silently ignore.
