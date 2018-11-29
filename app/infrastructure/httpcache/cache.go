package httpcache

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"
)

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
	sha := sha512.New()
	for key, values := range req.Header {
		sha.Write([]byte(key))
		for _, value := range values {
			sha.Write([]byte(value))
		}
	}
	sha.Write([]byte(req.Method))
	h := sha.Sum([]byte(req.URL.String()))
	return CacheKey(hex.EncodeToString(h))
}
