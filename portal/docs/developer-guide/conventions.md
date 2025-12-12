---
sidebar_label: Project Conventions
sidebar_position: 4
---

# Conventions Documentation

This document outlines the conventions used in the Digital Bank project for Git branching, commit messages, package naming, diagram creation, and logging practices. Adhering to these conventions ensures consistency and clarity across the codebase and documentation. Any questions or suggestions for improvements to these conventions can be raised in the meeting.

## Git

### Branching Conventions

#### During initial development phase:

- In the early stages of the development, the project would change repidly due to frequent requirement changes and as the developer knowledge/experience growth. Therefore, we will use a simpler branching strategy:
  - `main`: The main production branch.
  - `develop`: The main development branch where all features are integrated.
  - `feature/<service>/<short-description>`: Feature branches for new features (e.g., `feature/jwt-auth`).
  - `bugfix/<service>/<short-description>`: Bugfix branches for fixing bugs (e.g., `bugfix/token-expiry`).
  - `hotfix/<service>/<short-description>`: Hotfix branches for urgent fixes in production (e.g., `hotfix/login-issue`).
  - `release/<version>`: Release branches for preparing new versions (e.g., `release/v1.2.0`).

#### After the project is stable we may use a more advanced branching strategy:

- Use the following branch naming conventions:
  - `main`: The main production branch.
  - `<service>-main`: The main branch for each microservice (e.g., `auth-service-main`).
  - `<service>-dev`: The development branch for each microservice (e.g., `auth-service-dev`).
  - `feature/<service>/<short-description>`: Feature branches for new features (e.g., `feature/auth-service/jwt-auth`).
  - `bugfix/<service>/<short-description>`: Bugfix branches for fixing bugs (e.g., `bugfix/auth-service/token-expiry`).
  - `hotfix/<service>/<short-description>`: Hotfix branches for urgent fixes in production (e.g., `hotfix/auth-service/login-issue`).
  - `release/<service>/<version>`: Release branches for preparing new versions (e.g., `release/auth-service/v1.2.0`).

### Commit Message Conventions

- Use the following format for commit messages:

  ```
  <type>(<scope>): <subject>

  <body>

  <footer>
  ```

- **Type**: The type of change being made. Common types include:
  - `feat`: A new feature
  - `fix`: A bug fix
  - `refactor`: Code refactoring without changing functionality
  - `chore`: Maintenance tasks that do not affect the application code. (build process, dependencies, etc.)
  - `style`: Code style changes (formatting, missing semicolons, etc.)
  - `docs`: Documentation changes
  - `test`: Adding or updating tests
- **Scope**: The area of the codebase affected (e.g., `auth-service`, `http`, `usecases`, etc.)
- **Subject**: A brief description of the change (max 50 characters).
- **Body**: A more detailed explanation of the change, if necessary (wrap at 72 characters).
- **Footer**: Any relevant issue references or breaking changes.
- Use the imperative mood in the subject line (e.g., "Add feature" instead of "Added feature" or "Adds feature").
- Limit the subject line to 50 characters and the body to 72 characters per line.
- Separate the subject from the body with a blank line.
- Reference issues and pull requests in the footer using keywords like "Closes #123" or "Fixes #456".
- For breaking changes, include a "BREAKING CHANGE:" section in the footer with a description of the change and its impact.
- Example commit message:

  ```
  feat(auth-service): implement JWT authentication middleware

  Add a new middleware to handle JWT authentication for incoming requests.
  This middleware verifies the token and extracts user information.

  Closes #42
  ```

## Diagram Conventions

- Use PlantUML for creating diagrams.
- Follow a consistent style for all diagrams, including colors, fonts, and shapes.
- Clearly label all components, containers, and relationships in the diagrams.
- Use descriptive names for components and containers that reflect their purpose.
- Include a brief description of each component or container in the diagram.

## Logging Conventions

- **Library:** Use a structured logging library (**Zerolog**) for consistent JSON formatting.
- **Context:** Always include relevant context. `trace_id` and `instance_id` are mandatory.
- **Timestamp:** Use ISO 8601 format for timestamps.
- **Levels:** Use appropriate log levels (`INFO`, `DEBUG`, `WARN`, `ERROR`, `FATAL`) based on severity.
- **Clarity:** Ensure messages are machine-parsable. Avoid vague human text like "it failed".
- **Security:** **NEVER** log sensitive information (e.g., passwords, API keys, JWT tokens, PII).
- **Format:** All logs in production must be JSON.

### Standard Log Message Structure

All services must output logs following this schema:

```json
{
  "timestamp": "2025-11-03T08:15:30Z",
  "level": "INFO",
  "service": "ledger-service",
  "instance_id": "ledger-5f4a3b7c9d",
  "trace_id": "abc123xyz",
  "user_id": "12345",
  "action": "transfer_request",
  "status": "success",
  "details": {
    "amount": 10.0,
    "asset": "USDT",
    "to_user_id": "67890"
  },
  "message": "Transfer request processed successfully"
}
```

- [Logging Guideline](/docs/developer-guide/logging) has detailed examples for different scenarios.
