package cache

import "golang.org/x/sync/syncmap"

var DataCache = syncmap.Map{}

func GetMemoryCache(key string) (interface{}, bool) {
	return DataCache.Load(key)
}

func SetMemoryCache(key string, data interface{}) {
	DataCache.Store(key, data)
}
