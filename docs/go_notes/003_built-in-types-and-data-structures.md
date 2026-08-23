# Value semantics

## Mental model

Everything is **passed by value** (on assignment, function call), including pointers (value of ptr is address).
> NOTE: "Go has value semantics" does **not** mean "everything is deeply copied."

```go
b := a                      // variables
st := Point{X: 1, Y: 2}     // structs
st2 := st
    st.X = 7
    st2.X = 9
```

Some Go types contain references internally. The value is copied, but both copies may still refer to shared underlying data:
> Slices, Maps, Channels, Functions, Interfaces, Pointers

```go
sl := []int{1, 2, 3}    // slice
sl2 := sl
sl2[0] = 100            // both backing array changed
```

---

# Zero values

Every variable automatically has a valid default bit-level ZERO value.

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

A well-designed type often makes its zero value useful (no constructor required to use variable):

```go
var buf bytes.Buffer
    buf.WriteString("hello")
var mu sync.Mutex
var wg sync.WaitGroup

var xs []int                // slice
    xs = append(xs, 1)

var m map[string]int        // map (can be only read)
    fmt.Println(m["foo"])   // 0
    m["foo"] = 1            // panic        ;must: m = make(map[string]int)
```

which:
- reduce constructor requirements
- make composition easier

But not every type can or should have a fully useful zero value.
Prefer useful zero values where practical; use constructors where initialization establishes meaningful invariants or dependencies.

Constructors remain useful when:

* invariants must be established
* dependencies are required
* defaults differ from zero values
* internal maps/channels must be initialized
* validation is required

---

# Pointers

```go
x := 10
var p *int  // nil pointer          ;p == nil
p = &x      // get address of x     ;*p == 10
*p = 20     // get value at addr    ;*p == 20, x == 20

0xAA { p --> 0xAC       // p contains address of x
0xAC { x --> 20         // x contains 20
```

A pointer is an address to another value, usefule when:
- need **shared mutation**
- want to avoid copying a large value
- need to express optionality
- assign mutating methods to types (pointer receivers)

```go
func incrementPtr(x *int) {
    (*x)++      // modifies the caller's value.
}
```

Go does not have C-style pointer arithmetic. You cannot arbitrarily move a pointer through memory.

Do not use pointers simply because “objects should be pointers.” Small immutable-ish values are better passed as value:
- `time.Time`
- IDs
- coordinates
- money values
- configuration values

---

# Bytes

`byte` alias for `uint8` - one 8-bit integer (not character!)
`[]byte` - bytes slice - represents mutable binary data: network data, files, cryptography, encodings, compression, and I/O.

```go
data := []byte("hello")
data := []byte{72, 101, 108, 108, 111}      // mutable (unlike strings)
fmt.Println(string(data)) // Hello
data[0] = 'h'

data := []byte(str)     // convert str->[]byte
str := string(data)     // convert []byte->str
```

---

# Runes

`rune` alias for `int32` - one **Unicode code point** (1-4 bytes)

```go
r := 'ą'
fmt.Printf("%c\n", r)

rs := []rune("żółw")        // string -> runes slice
```

**BUT** runes are not necessarily visible characters - some consist multiple **code points**: `letter + combining accent`

```text
byte != rune != necessarily visible character
```

This matters when implementing cursor positions, text limits, or Unicode-aware UI behavior.

---

# Strings (stdlib package)

Immutable sequence of bytes: `pointer to bytes + length` (not array of characters); exact representation is an implementation detail.

```go
s := "żółw"
len(s)                  // counts **bytes**, not Unicode characters.
for i, r := range s {   // iterates UTF-8 decoded runes (int32 representing Unicode code point)
    fmt.Println(i, r)       // i - byte offset (not char number), r - rune
s[0]                    // returns byte (not whole UTF-8 character)
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
```
slice
 ├── pointer ─────► backing array
 ├── length
 └── capacity
```
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

# I/O and Streams

## io.Reader

`io.Reader`     - Something from which bytes can be read incrementally, consumer does not know from where

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

Files, HTTP response bodies, gzip decompressors, strings, network connections, and buffers can all act as readers:
```go
os.File
http.Response.Body
bytes.Buffer
strings.Reader
gzip.Reader
```

Each layer streams data through the next:
```text
HTTP body
   │
   ▼
gzip.Reader
   │
   ▼
json.Decoder
```

Bytes and EOF may occur together, process the bytes before acting on the error:
```go
n > 0
err == io.EOF       // sentinel error
```

## bufio

`bufio` adds buffering around readers/writers to reduce expensive underlying operations:
```go
reader := bufio.NewReader(file)
writer := bufio.NewWriter(file)
defer writer.Flush()                // remember - flush writers
```

`bufio.Scanner` is convenient for token/line-based processing:
```go
scanner := bufio.NewScanner(r)

for scanner.Scan() {
    line := scanner.Text()
}

if err := scanner.Err(); err != nil {
    return err
}
```
But Scanner has a token size limit by default. Large records require configuring its buffer or using another approach.

## Streaming versus buffering
```text
Buffering:
    input → load everything into memory → process   // e.g. 5 GB file
Streaming:
    input → process chunks progressively            // e.g. 64 KB chunks of 5 GB file

data, err := io.ReadAll(r)                          // buffers entire stream
io.Copy(dst, src)                                   // streams through bounded buffers
```

## String streams
```go
stream := strings.NewReader("string")     // strings.Reader implementing various io interfaces
```
