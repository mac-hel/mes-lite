# Observability

Three primary signals:

```text
logs
metrics
traces
```

Prefer structured logging:

```go
logger.Info("user created",
    "user_id", user.ID,
)
```

Propagate trace/request metadata through context rather than global mutable state.
