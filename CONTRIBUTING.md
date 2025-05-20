# Contributing to Awake

Thank you for considering contributing to Awake! This document explains how to contribute to the project, including our commit message convention that drives our semantic versioning system.

## Development Workflow

Please see [DEVELOPMENT.md](DEVELOPMENT.md) for detailed instructions on setting up your development environment and workflow.

## Commit Message Convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/) to automatically determine semantic version bumps and generate changelogs. Please adhere to this specification when creating commit messages.

The commit message should be structured as follows:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

The type must be one of the following:

- **feat**: A new feature (corresponds to a `minor` version bump)
- **fix**: A bug fix (corresponds to a `patch` version bump)
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to our CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files

### Optional Scopes

You can add an optional scope in parentheses to specify which part of the codebase is affected:

- **darwin**: For macOS-specific features or fixes
- **cli**: For command-line interface changes
- **docs**: For documentation changes
- **test**: For test-related changes

### Breaking Changes

Commits that introduce breaking changes should:
- Include `BREAKING CHANGE:` in the footer
- Or append a `!` after the type/scope

Breaking changes correspond to a `major` version bump.

### Examples

```
feat(darwin): add background mode option

Added a new command-line flag "-b" to run awake in background mode.
```

```
fix: correct duration parsing for time values
```

```
feat!: rename command-line arguments

BREAKING CHANGE: The "-t" flag has been renamed to "-d" for duration.
```

```
docs: update README with better installation instructions
```

## Pull Request Process

1. Ensure your code passes all tests and adheres to the project's coding standards
2. Update the README.md or other documentation if necessary
3. Verify that your changes have appropriate tests
4. Submit your pull request with a clear title and description

Thank you for contributing to Awake!
