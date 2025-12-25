# Contributing to Scout

Thank you for your interest in contributing to Scout!

## Commit Message Convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/) to automate versioning and changelog generation.

### Commit Message Format

```bash
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description | Version Bump |
| ------ | ------------- | -------------- |
| `feat` | A new feature | Minor |
| `fix` | A bug fix | Patch |
| `docs` | Documentation only changes | Patch |
| `style` | Code style changes (formatting, etc.) | Patch |
| `refactor` | Code refactoring | Patch |
| `test` | Adding or updating tests | Patch |
| `chore` | Maintenance tasks | Patch |
| `ci` | CI/CD changes | Patch |
| `build` | Build system changes | Patch |

### Breaking Changes

For breaking changes, add `!` after the type or include `BREAKING CHANGE:` in the footer:

```bash
feat!: remove deprecated --token flag

BREAKING CHANGE: The --token flag has been removed. Use environment variables instead.
```

Breaking changes trigger a **major** version bump.

### Examples

```bash
# Feature (bumps minor version)
git commit -m "feat(parser): add support for Gradle build files"

# Bug fix (bumps patch version)
git commit -m "fix(npm): handle missing package-lock.json gracefully"

# Documentation (bumps patch version)
git commit -m "docs: update installation instructions"

# Breaking change (bumps major version)
git commit -m "feat!: change output format for JSON exporter"
```

## Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes with conventional commits
4. Push to your fork: `git push origin feat/my-feature`
5. Open a Pull Request

## Release Process

Releases are automated via GitHub Actions:

1. When a PR is merged to `main`, the release workflow runs
2. Version is automatically bumped based on commit types:
   - `feat` commits → Minor version bump
   - `fix`, `docs`, etc. → Patch version bump
   - Breaking changes → Major version bump
3. CHANGELOG.md is updated automatically
4. A new GitHub release is created with binaries

## Pull Request Guidelines

Since every merge to `main` triggers an automatic release, we maintain high standards for what gets merged:

- **Only meaningful changes are merged** — PRs must provide clear value (new features, bug fixes, documentation improvements)
- **Use proper conventional commit format** — Your commit messages directly become the changelog entries
- **Keep PRs focused** — One feature or fix per PR makes reviews easier and changelogs cleaner
- **Write descriptive commit messages** — These are visible to all users in the release notes

PRs that don't meet these standards or lack a clear purpose will be requested for changes or closed.
