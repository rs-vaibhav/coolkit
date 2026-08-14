# Contributing to CoolKit

First off, thank you for considering contributing to CoolKit! It's people like you that make CoolKit such a great tool for student clubs. We welcome contributions of all kinds: bug fixes, new features, documentation improvements, and more.

## Code of Conduct

By participating in this project, you are expected to uphold our [Code of Conduct](CODE_OF_CONDUCT.md). Please report unacceptable behavior to the project maintainers.

## Getting Started

1. **Fork** the repository on GitHub.
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/your-username/coolkit.git
   cd coolkit
   ```
3. **Setup** the project. You have two options:
   - **Docker**: `make docker-up`
   - **Manual**: Ensure Go 1.22+ and PostgreSQL are installed. Copy `.env.example` to `.env`, update credentials, and run `make migrate-up`.

## Development Workflow

1. **Create a branch** for your work:
   - Feature: `feature/your-feature-name`
   - Bugfix: `fix/issue-description`
   - Docs: `docs/what-you-documented`
   ```bash
   git checkout -b feature/awesome-new-thing
   ```
2. **Make your changes**.
3. **Run checks**:
   ```bash
   make fmt
   make lint
   make test
   ```
4. **Commit** your changes using Conventional Commits.
5. **Push** to your fork and open a Pull Request (PR).

## Commit Message Convention

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

`type(scope): description`

Types:
- `feat:` A new feature
- `fix:` A bug fix
- `docs:` Documentation only changes
- `style:` Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- `refactor:` A code change that neither fixes a bug nor adds a feature
- `perf:` A code change that improves performance
- `test:` Adding missing tests or correcting existing tests
- `chore:` Changes to the build process or auxiliary tools and libraries such as documentation generation

Example: `feat(auth): add password reset functionality`

## PR Process

1. Provide a clear description of the changes in your PR.
2. Link any relevant issues (e.g., "Closes #123").
3. Ensure all CI checks pass.
4. Request a review from one of the maintainers.
5. Make any requested changes.

## Coding Standards

- Run `go fmt ./...` or `make fmt` before committing.
- Ensure `golangci-lint` passes (`make lint`).
- Write tests for new features or bug fixes.
- Keep functions small and focused on a single responsibility.
- Document exported types, functions, and methods.

## Adding a New Feature (Walkthrough)

1. **Model**: Add or update the GORM model in `internal/model`.
2. **Repository**: Add database access methods in `internal/repository`. Create an interface and a concrete implementation.
3. **Service**: Add business logic in `internal/service`. Call the repository methods here.
4. **Handler**: Add HTTP request/response parsing and validation in `internal/handler`.
5. **Route**: Register the new handler in `cmd/server/main.go` (or wherever your router is defined).
6. **Tests**: Add unit tests for your service and handler.

## First-Time Contributors

Look for issues labeled **"good first issue"** or **"help wanted"**. These are typically easier tasks like adding small tests, fixing minor bugs, or improving documentation.

## Project Architecture

CoolKit uses a layered architecture to separate concerns. Please review our [Architecture Documentation](docs/architecture.md) before making significant changes.

## Getting Help

If you need help or have questions, please:
- Open an issue for bug reports or feature requests.
- Start a discussion in the GitHub Discussions tab for general questions or ideas.
