---
sidebar_label: Project Structure
sidebar_position: 2
---


# Project Structure

```plaintext
bitka/
├── .env                       # Environment variables for local dev
├── .env.template              # Template for onboarding devs
├── .gitignore
├── .vscode/                   # Editor settings (extensions, formatter)
│   └── settings.json
├── bitka.code-workspace       # VSCode workspace to link multiple repositorys
├── docker-compose.yml         # Local environment stack (db, kafka, etc.)
├── LICENSE
├── Makefile                   # build, test, lint, run commands
├── README.md
├── go.work                    # Go workspace tying multiple modules together
│
├── docs/
│   ├── architecture/                    # High-level system design
│   │   ├── system-design.md             # Big picture: services, data flow, context diagram
│   │   ├── database-schema.md           # ERD + explanation of major tables
│   │   └── diagrams/                    # All PlantUML diagrams (.puml only)
│   │
│   └── api/                             # API documentation (OpenAPI + AsyncAPI)
│       ├── openapi.yaml                 # Main entrypoint (references components + paths)
│       ├── asyncapi.yaml                # (optional) Kafka event contract
│       │
│       ├── components/                  # Reusable parts (schemas, responses)
│       │   ├── schemas/
│       │   │   ├── Order.yaml           # Schema: how an Order looks
│       │   │   └── User.yaml
│       │   ├── responses/
│       │   │   └── Error.yaml           # Standard error format
│       │   └── parameters/              # Query/path parameters (optional)
│       │       └── pagination.yaml
│       │
│       └── paths/                       # Endpoint definitions
│           ├── orders.yaml              # /orders
│           ├── auth.yaml                # /auth (login, signup)
│           ├── wallet.yaml              # /wallet (balance, transfers)
│           └── (other paths)
│
├── pkg/                       # Shared, reusable libraries
│   ├── go.mod                 # Module bitka/pkg
│   ├── config/                # Shared config loader (env → struct) and helpers
│   ├── database/              # DB utilities (gorm, pgx, etc.)
│   ├── logger/                # Zerolog setup (standard logger)
│   ├── middleware/            # Auth middleware usable by services
│   ├── response/              # Standard API response wrapper
│   └── token/                 # JWT & JWX utilities
│       ├── jwx_manager.go
│       ├── model.go
│       └── validator.go
│
├── services/                  # All microservices
│   ├── Dockerfile.auth        # Dockerfile for auth service
│   ├── Dockerfile.(other)     # (other services)
│   ├── auth/                  # (Example shown; other services follow same structure)
│   │   ├── go.mod             # module bitka/services/auth
│   │   ├── cmd/server/
│   │   │   └── main.go        # service entrypoint
│   │   └── internal/          # Clean Architecture layout
│   │       ├── app/           # App bootstrap (server creation, DI)
│   │       │   └── server.go
│   │       ├── domain/        # Business entities + interfaces
│   │       ├── repository/    # Remote source implementations
│   │       ├── usecase/       # Application logic (interactors)
│   │       └── delivery/      # HTTP DTO, routing, handlers / GRPC / events
│   │           └── http/
│   │               └── handler.go
│   │
│   ├── account/
│   ├── outbox/                # Outbox pattern service
│   └── (other services planned: ledger, order, matcher, marketdata, deposit, withdraw, notify...)
│
└── tools/ (optional in future) # Scripts, linters, generators
    ├── lint.sh                 # Run all linters
    └── generate.sh             # Code generation scripts
```