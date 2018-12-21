package httpcache

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"net/http"
)

// CacheKey represents a key in base64 () of sha512 (64 bytes).
type CacheKey string

type CacheValue []byte

type Cache interface {
	// Get returns the value or nil if the key is not found.
	Get(CacheKey) (CacheValue, error)
	// Set stores the key and value.
	Set(CacheKey, CacheValue) error
	// Delete deletes the key. Do not return error even if the key is not found.
	Delete(CacheKey) error
}

func computeCacheKey(req *http.Request) CacheKey {
	var b bytes.Buffer
	for key, values := range req.Header {
		b.Write([]byte(key))
		for _, value := range values {
			b.Write([]byte(value))
		}
	}
	b.Write([]byte(req.Method))
	b.Write([]byte(req.URL.String()))
	h := sha512.Sum512(b.Bytes())
	e := base64.StdEncoding.EncodeToString(h[:])
	return CacheKey(e)
}
