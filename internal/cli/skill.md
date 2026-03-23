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
- The canonical body sections are `Objective`, `Context`, `Research`, `Plan`, `Outputs`, `Verification`, `Review`, and `Notes`.

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
