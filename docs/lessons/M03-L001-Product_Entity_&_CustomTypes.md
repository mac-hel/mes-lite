## Lesson 3.1 — Custom Types & iota: Product Entity

**Business Context:** The company manufactures ventilation components. Products fall into categories (ventilation, filters, ducts, mounting hardware, other). Each product has a SKU, name, category, unit of measure, and active status.

**Why custom types + iota instead of enums?** Go has no `enum` keyword. Instead, idiomatic Go uses a custom integer type + `iota` constant generator. This gives you compile-time type safety — you can't accidentally pass a raw `int` where a `ProductCategory` is expected, unlike PHP where enums need a backing library.

**Why Stringer?** The `fmt.Stringer` interface (`String() string`) is Go's equivalent of `__toString()`. Implementing it means `fmt.Print`, `%s`, `%v` all produce a readable label — useful for logging and JSON serialization.

**Go concepts introduced:**
- **custom types** (`type ProductCategory int`) — give primitive types domain meaning and compile-time type safety
- **iota** — Go's constant generator; each const in a block auto-increments, the idiomatic replacement for enums
- **Stringer** (`fmt.Stringer` interface) — `String() string` gives types a human-readable representation; used by `fmt`, `%s`, `%v`, and logging
- **value semantics** — `Product` is passed/returned by value (like `Employee`), since it's a small struct with no shared mutable state

**Idioms demonstrated:**
- Custom type + iota instead of enum keyword (Go has no enums)
- Stringer implementation for readable output
- Zero value of `ProductCategory` = `CategoryVentilation` (0), which is a meaningful default

**Common mistakes to avoid:**
- Using string constants instead of custom int types (loses type safety: `"Ventilation"` is just a string, any string fits)
- Forgetting the `default` case in String() switch (future values would print nothing)
- Exporting iota values without a comment block on the const

**Exercises:**
- Add a new category `CategoryElectrical` between `CategoryVentilation` and `CategoryFilter` — what happens to the numbering? How would you fix it?
- Implement `String()` without a `default` case. What does `fmt.Println(CategoryVentilation + 99)` print?

**Interview questions:**
- *Why doesn't Go have enums?* — Go prefers simplicity. iota + custom type gives 80% of the value (type-safe constants) without dedicated syntax. True enums with exhaustive switch checking can be approximated with linters.
- *Why custom types over string constants?* — `type Status string` prevents accidentally passing a raw string where a status is expected. The compiler enforces the distinction.
- *When would you use a pointer receiver for String()?* — Almost never. String() is a read operation; value receivers are safe, simpler, and work even on nil (a pointer receiver on nil would panic).
