package test

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestOne(t *testing.T) {
	s := "12345678"

	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	fmt.Printf("origin: %s, sha256 hash: %x\n", s, bs)

	ss := fmt.Sprintf("%x", bs)
	h = sha256.New()
	h.Write([]byte(ss))
	bs = h.Sum(nil)
	fmt.Printf("origin: %s, sha256 hash: %x\n", s, bs)

	ss = fmt.Sprintf("%x", bs)
	h = sha256.New()
	h.Write([]byte(ss))
	bs = h.Sum(nil)
	fmt.Printf("origin: %s, sha256 hash: %x\n", s, bs)

	ss = fmt.Sprintf("%x", bs)
	h = sha256.New()
	h.Write([]byte(ss))
	bs = h.Sum(nil)
	fmt.Printf("origin: %s, sha256 hash: %x\n", s, bs)

	ss = fmt.Sprintf("%x", bs)
	h = sha256.New()
	h.Write([]byte(ss))
	bs = h.Sum(nil)
	fmt.Printf("origin: %s, sha256 hash: %x\n", s, bs)

	ss = fmt.Sprintf("%x", bs)
	h = sha256.New()
	h.Write([]byte(ss))
	bs = h.Sum(nil)
	fmt.Printf("origin: %s, sha256 hash: %x\n", s, bs)

	ss = fmt.Sprintf("%x", bs)
	h = sha256.New()
	h.Write([]byte(ss))
	bs = h.Sum(nil)
	fmt.Printf("origin: %s, sha256 hash: %x\n", s, bs)

	ss = fmt.Sprintf("%x", bs)
	h = sha256.New()
	h.Write([]byte(ss))
	bs = h.Sum(nil)
	fmt.Printf("origin: %s, sha256 hash: %x\n", s, bs)
}

// password: 12345678
// sha256 EIGHT times result: 04a24e8195382cbfe6c81dda873d2be49b13c1bd09b01f0bfeeba952de3c59cd
