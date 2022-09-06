package dns256

import (
	"encoding/json"
	"time"
)

type Cache struct{}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Get(x string, now time.Time) *CacheEntry {
	return nil
}

func (c *Cache) Put(ent CacheEntry) {

}

type CacheEntry struct {
	Key       string
	Value     json.RawMessage
	ExpiresAt time.Time
}
