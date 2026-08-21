# Arrays

Array length is part of its type:

```go
[3]int
[4]int
```

are different types.

Arrays have value semantics:

```go
b := a
```

copies the array.

Passing an array to a function copies it unless a pointer is used.

Arrays are less common directly; slices are usually the user-facing abstraction.

---

# Slices

A slice is a small descriptor referring to an underlying array, has current length and max capacity.
`pointer → backing array + length + capacity`
Lookup: O(n)

```go
var s []int                 // nil,     len(s) == 0 , cap(s) == ?
s := []int{1, 2, 3}         // non-nil, len(s) == 3 , cap(s) == ?
s := make([]int, 4, 10)     // non-nil, len(s) == 4 , cap(s) == 10

b := s                  // copies reference; NOT independent copy: both refer to the dame backing array
b[0] = 99               // now s[0] == 99

s = append(s, 23)       // if capacity insufficient, append will: allocate new backing array, copy elements
slices.Equals(s1, s2)
```

## Subslice aliases memory
Slice and subslice share backing storage - changes may be visible through both slices.
```go
huge := loadHugeBuffer()
tiny := huge[2:10]              // elements 2-9; [:10] elements 0-9; [2:] elements 2-LAST
```

Copy when necessary:
```go
tiny := append([]byte(nil), huge[:10]...)
    // or
tiny := make([]byte, 10)
copy(tiny, huge[:10])
```

## Nil vs empty slice
Both are valid for operations:  `range, append, len, cap`
But they can differ under:      `== nil`, JSON encoding, reflection, some API contracts
```go
var a []int        // nil,      length 0
b := []int{}       // non-nil,  length 0
c := make([]int,0) // non-nil,  length 0
```

## sorting
    import "sort"           // package for sorting Slices and user-defined Collections
	sort.Slice(prods, func(i, j int) bool {
		a, b := prods[i], prods[j]
		switch sortKey {
		case "-sku":
			return a.SKU > b.SKU
		case "name":
			return a.Name < b.Name
		case "-name":
			return a.Name > b.Name
		default:
			return a.SKU < b.SKU
		}
	})

---

# Maps

Runtime-managed structures with reference-like behavior.
Lookup: O(1)

```go
m := make(map[string]User)                  // initialize (memory) - map must be initialized before writing
m := make(map[User]struct{}, len(users))    // Map Set pattern
m["x"] = 1                                  // write entry

m2 := m                     // copies reference; NOT independent copy: mutation of one's entries is visible through the other
```

## Nil map
Safe to use (but 0 values) except for writing.
```go
var m map[string]int    // nil map
v := m["x"]             // 0
len(m)                  // 0
delete(m,"x")           // safe
m["x"] = 1              // INVALID - panic
```

## Missing key
Zero is returned for absent values; use `ok`.
```go
v := m[key]         // returns the value type's zero value when absent
v, ok := m[key]     // to distinguish absent from a stored zero value
    if !ok { ... }
```

## Map iteration
Sort keys explicitly when deterministic output matters.
```go
for k := range m    // iteration order is unspecified
```

## Map concurrency
Ordinary maps are not generally safe for unsynchronized concurrent read/write access, typically needs:
```go
sync.Mutex
sync.RWMutex
sync.Map        // specialized; a normal map protected by a mutex is often clearer
```

## `comparable` - in other docs

---

# Strings (stdlib package)

Immutable byte data: `pointer to bytes + length`; exact representation is an implementation detail.

```go
len(s)                  // counts **bytes**, not Unicode characters.
for _, r := range s     // iterates UTF-8 decoded runes (int32 representing Unicode code point)
```
strings.Contains
strings.ToLower
strings.TrimSpace
strings.HasPrefix(s, prefix)
strings.TrimPrefix(s, prefix)
strings.Join([]string{"str1", "str2"}, "\n"))
strings.Count(str, substr)              // count non-overlapping substr in str

strconv.Atoi(s)/ParseInt(S,10,0)        // str -> int
strconv.ParseBool(b)                    // str -> bool

## Stringer interface
`fmt.Stringer` interface (`String() string`) - `fmt.Print`, `%s`, `%v` all produce a readable label; useful for logging and JSON serialization.

```go
func (c ProductCategory) String() string { switch c { case CategoryVentilation: return "Ventilation"
```

---

# Streams `import io`
- io.Reader     // interface
- io.EOF        // sentinel error
- stream := strings.NewReader("string")     // strings.Reader implementing various io interfaces

---

# Zero values

Every variable has a zero value.

Examples:

```go
int       -> 0
bool      -> false
string    -> ""
pointer   -> nil
slice     -> nil
map       -> nil
channel   -> nil
interface -> nil
```

A well-designed type often makes its zero value useful:

```go
var mu sync.Mutex
var buf bytes.Buffer
```

But not every type can or should have a fully useful zero value.
Prefer useful zero values where practical; use constructors where initialization establishes meaningful invariants or dependencies.

Constructors remain useful when:

* invariants must be established
* dependencies are required
* defaults differ from zero values
* internal maps/channels must be initialized
* validation is required
