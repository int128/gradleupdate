// Package memcache is
// forked from https://github.com/gregjones/httpcache/blob/master/memcache/appengine.go
package memcache

import (
	"context"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

// Cache is an implementation of httpcache.Cache that caches responses in App
// Engine's memcache.
type Cache struct {
	Context context.Context
}

// cacheKey modifies an httpcache key for use in memcache.  Specifically, it
// prefixes keys to avoid collision with other data stored in memcache.
func cacheKey(key string) string {
	return "httpcache:" + key
}

// Get returns the response corresponding to key if present.
func (c *Cache) Get(key string) (resp []byte, ok bool) {
	item, err := memcache.Get(c.Context, cacheKey(key))
	if err != nil {
		if err != memcache.ErrCacheMiss {
			log.Errorf(c.Context, "error getting cached response: %v", err)
		}
		return nil, false
	}
	return item.Value, true
}

// Set saves a response to the cache as key.
func (c *Cache) Set(key string, resp []byte) {
	item := &memcache.Item{
		Key:   cacheKey(key),
		Value: resp,
	}
	if err := memcache.Set(c.Context, item); err != nil {
		log.Errorf(c.Context, "error caching response: %v", err)
	}
}

// Delete removes the response with key from the cache.
func (c *Cache) Delete(key string) {
	if err := memcache.Delete(c.Context, cacheKey(key)); err != nil {
		log.Errorf(c.Context, "error deleting cached response: %v", err)
	}
}

// New returns a new Cache for the given context.
func New(ctx context.Context) *Cache {
	return &Cache{ctx}
}
