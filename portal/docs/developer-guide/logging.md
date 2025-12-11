---
sidebar_label: Logging Guideline
sidebar_position: 4
---

# How to Log

The "Golden Rule" of logging in Bitka is: **Never create a new logger. Always extract it from the Context.**

## 🔑 The Pattern

Every function in the data flow (Handler $\rightarrow$ Usecase $\rightarrow$ Repository) accepts `context.Context`. This context carries the **Logger** and the **Trace ID**.

```go
// ❌ BAD: Uses global logger (No Trace ID, No Context)
log.Info().Msg("User created")

// ✅ GOOD: Extracts context-aware logger
logger.From(ctx).Info().
    Str("action", "user_create").
    Str("status", "success").
    Msg("User created")
```

## 🌊 Data Flow 1: HTTP Request (Fiber)

In Fiber, the middleware automatically injects the logger into the `UserContext`.

### **1. Delivery Layer (Handler)**

```go
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    // 1. Get Context (Contains Trace ID & Logger)
    ctx := c.UserContext()

    // 2. Log the Action
    logger.From(ctx).Info().
        Str("action", "login_attempt").
        Str("email", req.Email).
        Str("status", "processing").
        Msg("Processing login")

    // 3. Pass Context down
    return h.usecase.Login(ctx, req.Email, req.Password)
}
```

### **2. Usecase Layer**

```go
func (u *AuthUsecase) Login(ctx context.Context, email, password string) error {
    // Logic...
    if err != nil {
        // Log Errors with Stack Trace support
        logger.From(ctx).Error().
            Err(err).
            Str("action", "login_logic").
            Str("status", "failed").
            Msg("Invalid credentials or internal error")
        return err
    }
    return nil
}
```

### **3. Repository Layer**

```go
func (r *Repo) FindUser(ctx context.Context, email string) (*User, error) {
    // GORM automatically uses the context logger if configured,
    // but if you need manual logging:
    logger.From(ctx).Debug().
        Str("action", "db_query").
        Str("query", "select_user_by_email").
        Msg("Executing query")

    // ... DB call
}
```

## 📨 Data Flow 2: Kafka Consumer

Kafka consumers are background workers. They don't have an HTTP request, so we must **generate** a context for them.

```go
// internal/delivery/kafka/handler.go

func (h *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        // 1. Extract or Generate Trace ID
        // Try to read "trace_id" from Kafka Headers, otherwise generate new
        traceID := extractTraceID(msg.Headers)

        // 2. Create Logger Context
        // We manually start the logging session for this message
        l := log.With().
            Str("trace_id", traceID).
            Str("topic", msg.Topic).
            Int64("offset", msg.Offset).
            Logger()

        // 3. Create Standard Context
        ctx := context.Background()
        ctx = logger.WithContext(ctx, &l)

        // 4. Call Usecase
        logger.From(ctx).Info().
            Str("action", "kafka_consume").
            Str("status", "start").
            Msg("Processing message")

        err := h.usecase.ProcessEvent(ctx, msg.Value)

        if err != nil {
            logger.From(ctx).Error().Err(err).Str("status", "failed").Msg("Processing failed")
        }
    }
    return nil
}
```

## 📤 Data Flow 3: Kafka Producer

When sending events, we want to log that we _sent_ it, preserving the Trace ID of the request that triggered it.

```go
// internal/repository/kafka/producer.go

func (p *Producer) PublishUserCreated(ctx context.Context, event UserCreatedEvent) error {
    // Use the logger from the context (inherits trace_id from HTTP request)
    l := logger.From(ctx)

    l.Info().
        Str("action", "publish_event").
        Str("topic", "user.events.v1").
        Str("event_type", "UserCreated").
        Interface("payload", event). // Careful not to log PII here!
        Msg("Emitting Kafka event")

    // Send logic...
}
```

## 🎚️ Log Level Guide

| Level     | When to use                                                               | Behavior                  | Example                                                               |
| :-------- | :------------------------------------------------------------------------ | :------------------------ | :-------------------------------------------------------------------- |
| **FATAL** | **Startup Only.** The application cannot start (Missing Config, DB down). | **Exits App (`os.Exit`)** | `{"level":"FATAL", "action":"init_db", "error":"connection refused"}` |
| **ERROR** | Operation failed, user gets 500. Action required by Ops/Devs.             | Continues running         | `{"level":"ERROR", "action":"db_query", "error":"timeout"}`           |
| **WARN**  | Unexpected state, but handled. Retries, deprecated API usage, 400 errors. | Continues running         | `{"level":"WARN", "action":"validate_input", "reason":"bad_email"}`   |
| **INFO**  | Key business events. Happy path tracking.                                 | Continues running         | `{"level":"INFO", "action":"order_placed", "order_id":"123"}`         |
| **DEBUG** | detailed payload dumps, logic flow.                                       | Hidden in Prod            | `{"level":"DEBUG", "action":"calc_fee", "vars":{...}}`                |

### ⚠️ The Danger of `FATAL`

Unlike `ERROR` (which just logs that something went wrong), `FATAL` usually performs a system call to **exit the program immediately** (`os.Exit(1)`).

- **In a Microservice:** If you call `log.Fatal()` inside an HTTP handler (e.g., failed to save a user), **it kills the entire server**. All other active requests will drop instantly. Kubernetes will see the crash and restart the pod.
- **In Main:** This is the only safe place to use it.

### Code Examples

#### ✅ GOOD Usage (Startup in `main.go`)

It is acceptable to kill the app here because it's useless without the database.

```go
func main() {
    // ... load config ...

    db, err := database.Connect(cfg)
    if err != nil {
        // FATAL: Log the error and kill the process immediately
        log.Fatal().
            Err(err).
            Str("action", "db_connect").
            Msg("Could not connect to database, shutting down")
    }
}
```

#### ❌ BAD Usage (Inside a Handler)

**Never** do this. This turns a simple failed login attempt into a server crash.

```go
func (u *AuthUsecase) Login(ctx context.Context, email string) error {
    user, err := u.repo.Find(ctx, email)
    if err != nil {
        // ☠️ DANGER: This kills the whole container!
        // logger.From(ctx).Fatal().Err(err).Msg("Database error")

        // ✅ CORRECT: Just log error and return
        logger.From(ctx).Error().Err(err).Msg("Database error")
        return err
    }
    return nil
}
```