package httpcache

import (
	"context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

type AppEngineMemcache struct {
	Context context.Context
}

func (c *AppEngineMemcache) Get(k CacheKey) (CacheValue, error) {
	item, err := memcache.Get(c.Context, string(k))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil
		}
		log.Errorf(c.Context, "Error while getting item from memcache: %s", err)
		return nil, err
	}
	return item.Value, nil
}

func (c *AppEngineMemcache) Set(k CacheKey, v CacheValue) error {
	err := memcache.Set(c.Context, &memcache.Item{
		Key:   string(k),
		Value: v,
	})
	if err != nil {
		log.Errorf(c.Context, "Error while setting item to memcache: %s", err)
	}
	return err
}

func (c *AppEngineMemcache) Delete(k CacheKey) error {
	err := memcache.Delete(c.Context, string(k))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil
		}
		log.Errorf(c.Context, "Error while deleting item in memcache: %s", err)
	}
	return err
}
