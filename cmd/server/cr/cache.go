package cr

import (
	"context"
	"sync"
	"time"

	"github.com/ReanSn0w/kincong/internal/utils"
)

func newCache(ctx context.Context) *Cache {
	cache := &Cache{
		data: make(map[string]cacheValue),
	}

	go cache.runCleaner(ctx)

	return cache
}

type Cache struct {
	mx   sync.RWMutex
	data map[string]cacheValue
}

type cacheValue struct {
	domainData *DomainInfo
	subnets    []utils.ResolvedSubnet
	createdAt  time.Time
}

func (c *Cache) GetSubnet(key utils.Value) ([]utils.ResolvedSubnet, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if val, ok := c.data["subnet_"+string(key)]; ok {
		if val.subnets == nil {
			return nil, false
		}

		return val.subnets, true
	}
	return nil, false
}

func (c *Cache) SetSubnet(key utils.Value, subnets []utils.ResolvedSubnet) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.data["subnet_"+string(key)] = cacheValue{
		subnets:   subnets,
		createdAt: time.Now(),
	}
}

func (c *Cache) GetDomainInfo(key utils.Value) (*DomainInfo, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if val, ok := c.data["domain_"+string(key)]; ok {
		if val.domainData == nil {
			return nil, false
		}

		return val.domainData, true
	}

	return nil, false
}

func (c *Cache) SetDomainInfo(key utils.Value, ipInfo DomainInfo) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.data["domain_"+string(key)] = cacheValue{
		domainData: &ipInfo,
		createdAt:  time.Now(),
	}
}

func (c *Cache) runCleaner(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	day := 24 * time.Hour

	for {
		select {
		case <-ticker.C:
			c.mx.Lock()
			for key, value := range c.data {
				if time.Since(value.createdAt) > day {
					delete(c.data, key)
				}
			}
			c.mx.Unlock()
		case <-ctx.Done():
			return
		}
	}
}
