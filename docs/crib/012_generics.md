# Generics

Use type parameters when the algorithm genuinely operates over a family of types:

```go
func Contains[T comparable](xs []T, x T) bool
```

Do not replace ordinary interfaces with generics automatically.

Rough distinction:

```text
interface -> runtime behavioral abstraction
generic   -> compile-time type abstraction
```
