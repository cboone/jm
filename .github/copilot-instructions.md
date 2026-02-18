# GitHub Copilot Instructions for FM

For full project conventions, see AGENTS.md in the repository root.

## PR Review

When reviewing pull requests, do not flag the following patterns as issues.
Each is an intentional project convention:

- **CLAUDE.md is a symlink**: CLAUDE.md is a symlink to AGENTS.md, not a duplicated file. The symlink target content appearing in diffs is expected behavior. Do not flag it as duplication.
