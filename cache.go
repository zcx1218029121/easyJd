package main

import "sync"

type Cache interface {
	set(key string, info *PageInfo)
	get(param string) (*PageInfo, bool)
}

type SimpleCache struct {
	sync.Map
}

func (c *SimpleCache) set(key string, info *PageInfo) {
	c.Map.Store(key, info)
}
func (c *SimpleCache) get(param string) (*PageInfo, bool) {
	if info, ok := c.Map.Load(param); ok {
		return info.(*PageInfo), true
	}
	return nil, false
}
