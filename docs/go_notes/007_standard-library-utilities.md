# Time `import time`
- var now time.Time = time.Now().UTC()
- var ttl time.Duration = 15 * time.Minute, now.Add(ttl)
- ttl.IsZero()
- s := now.Format(time.RFC3339) / now, err := time.Parse(time.RFC3339, s)

# Rand / UUID
- `crypto/rand` and `encoding/hex` for standard-library UUID-shaped IDs
