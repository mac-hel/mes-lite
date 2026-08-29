// Package ids generates the UUID-shaped identifiers used by business slices.
//
// The identifiers are random (version 4) UUIDs built from the standard library
// only. Business slices name their own identifier fields; they do not need
// their own generator.
package ids

import (
	"crypto/rand"
	"encoding/hex"
)

// New returns a random version 4 UUID in the canonical 8-4-4-4-12 form.
//
// New does not return an error. Since Go 1.24 crypto/rand.Read never fails: it
// fills the buffer completely or crashes the program irrecoverably. A machine
// whose randomness source is unusable cannot serve requests at all, so callers
// get a value instead of an error branch that can never be taken.
func New() string {
	var b [16]byte

	// The error is intentionally discarded: crypto/rand.Read is documented to
	// never return one, and to crash the program if the OS source fails.
	_, _ = rand.Read(b[:])

	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant

	var dst [36]byte
	hex.Encode(dst[0:8], b[0:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], b[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], b[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], b[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:36], b[10:16])

	return string(dst[:])
}
