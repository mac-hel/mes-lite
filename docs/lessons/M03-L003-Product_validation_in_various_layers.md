Lesson 3.3 is complete and `ROADMAP.md` is updated.

**Implemented**
- Added product domain validation in `internal/products/product.go`.
- Added category validity via `ProductCategory.Valid()`.
- Normalized product inputs with `strings.TrimSpace`.
- Added `Product.UpdateDetails()` for update normalization.
- Added `ErrInvalidProduct` with wrapped validation errors.
- Enforced validation in handlers and in-memory store.
- Added request validation tags for category range `0..4`, keeping OpenAPI/Fuego validation aligned.
- Added table-driven tests for validation and category behavior.

**Roadmap**
- Marked Milestone 3 as completed.
- Marked Lesson 3.3 as completed.
- Advanced current milestone to `4 - Production Registration`.
- Updated Knowledge Matrix: `Table Tests` checked.
- Bumped roadmap version to `1.3`.

**Verification**
- `go test ./... -count=1` passes.
- `make lint` passes with `0 issues`.

**Lesson Summary**
- Business rule: invalid products should not enter the system from HTTP or direct store usage.
- Go concepts: value semantics, pointer receiver for mutation, sentinel errors, error wrapping, table-driven tests.
- Idioms: keep validation close to the domain model, keep HTTP DTO validation separate from domain invariants, concrete first.
- Common mistake: relying only on HTTP validation tags. Domain rules must also protect non-HTTP callers.

**Code Review**
No blocking issues found. A stricter `SKU` value object type could be introduced later, but delaying it is the better Go choice for now because `string` still keeps the API simple and the current validation solves the real problem.
