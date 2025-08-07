package cache

import (
	"sync"
)

// TxNonce structure for a pair sender and nonce
type TxNonce struct {
	Address string
	Nonce   uint64
}

// AnteCache is a cache used by CheckAndSetEthSenderNonce to check that a specific TxNonce exists in a mempool
// TODO to be removed once it is correctly implemented in cosmos sdk
type AnteCache struct {
	mu    sync.RWMutex
	cache map[TxNonce]bool
	size  int
}

func NewAnteCache() *AnteCache {
	return &AnteCache{
		cache: make(map[TxNonce]bool),
	}
}

// Set the TxNonce
func (c *AnteCache) Set(address string, nonce uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := TxNonce{address, nonce}
	c.cache[key] = true
	c.size++
}

// Delete the TxNonce
func (c *AnteCache) Delete(address string, nonce uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := TxNonce{address, nonce}
	delete(c.cache, key)
}

// Exists check if the TxNonce existz
func (c *AnteCache) Exists(address string, nonce uint64) bool {
	key := TxNonce{address, nonce}
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.cache[key]
	return ok
}

func (c *AnteCache) Size() int {
	return c.size
}
