# Contributing to Bolte Bridge

This guide explains how to open issues and pull requests so your
contributions can be reviewed and merged smoothly!

## Before you start

- Search [existing issues](https://github.com/kcexn/bolte-bridge/issues) and
  [pull requests](https://github.com/kcexn/bolte-bridge/pulls) to avoid
  duplicating work that's already reported or in progress.
- For anything more than a trivial fix, open an issue first to discuss the
  change. This saves you from investing effort in an approach that may not be
  accepted.

## Opening an issue

Good issues are easy to act on. When reporting a **bug**, please include:

- A clear, descriptive title.
- What you expected to happen and what actually happened.
- Steps to reproduce, ideally a minimal example.
- Relevant environment details (OS, Go version, Bolte Bridge version/commit)
  and any log output.

When requesting a **feature**, describe the problem you're trying to solve and
why, not just the solution you have in mind. That context helps us find the
best fit for the project.

## Opening a pull request

1. **Fork** the repository and create a branch off `main` with a descriptive
   name (e.g. `fix/smtp-timeout` or `feat/matrix-threading`).
2. **Make focused changes.** Keep each pull request scoped to a single concern;
   smaller PRs are easier to review and merge.
3. **Follow the existing code style.** Run `gofmt`/`golangci-lint` and make sure
   the linter is happy.
4. **Add or update tests** for your change, and make sure the full test suite
   (`go test ./...`) passes locally.
5. **Write clear commit messages** that explain what changed and why.
6. **Reference related issues** in the PR description (e.g. `Closes #123`).
7. **Fill out the PR description**: summarize what the change does, how you
   tested it, and any trade-offs or follow-up work.

Once opened, CI will run automatically. Please address review feedback and
keep your branch up to date with `main`.
