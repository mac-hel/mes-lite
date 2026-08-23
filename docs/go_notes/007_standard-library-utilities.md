# Time

* store time in UTC format: `now.UTC()`
* convert to local zones on presentation boundaries

## Types

Go's primary time types are:
```go
import time

time.Time           // point in time
time.Duration       // elapsed time
```

## Operations

```go
var ttl time.Duration = 15 * time.Minute        // NOT `= 900` because no info about units

var now time.Time = time.Now().UTC()
now.Add(ttl)
myt.Equal(now)

ttl.IsZero()
```

## Formatting / Parsing

Go's unusual formatting/parsing layout uses this reference timestamp: `Mon Jan 2 15:04:05 MST 2006`

```go
            \/format      \/value
time.Parse("2006-01-02", "2026-08-22")

s := now.Format(time.RFC3339)
err := time.Parse(time.RFC3339, s)
```

---

# 7. Encoding

## JSON
`encoding/json` converts between Go values and JSON.
```go
type UserDTO struct {
    ID   string     `json:"id"`
    Groups []string `json:"groups"` // nil -> null, empty slice -> []
    age  int                        // priv field ommited
}

data, err := json.Marshal(user)         // encoding (to string)
err := json.Unmarshal(data, &user)      // decoding (back to struct)

var x any
json.Unmarshal(data, &x)                // no types (ambiguity), e.g. number becomes float64
```

For HTTP streaming, use:
```go
decoder := json.NewDecoder(r)
decoder.DisallowUnknownFields()     // unknown fields ignored
err := decoder.Decode(&dto)
    // instead
body, _ := io.ReadAll(r)
json.Unmarshal(body, &dto)
```

---

## CSV
`encoding/csv` - CSV handling works naturally as a stream of records.

Reader:
```go
reader := csv.NewReader(r)
for {
    record, err := reader.Read()    // record []string

    if errors.Is(err, io.EOF) {
        break
    }
    if err != nil {
        return err
    }

    process(record)
}

reader.Comma = ';'      // set delimiter
```

Writer:
```go
writer := csv.NewWriter(w)
defer writer.Flush()

if err := writer.Error(); err != nil {
    return err
}
```

CSV parsing is usually an infrastructure/boundary concern. Convert CSV records into DTOs/domain commands rather than allowing raw column positions to leak through the application.

---

# Rand / UUID
- `crypto/rand` and `encoding/hex` for standard-library UUID-shaped IDs
