package ccache

import (
	"bm/internal/tool"
	"github.com/coocood/freecache"
	"runtime/debug"
	"sync"
)

const DefaultSize = 100 * 1024 * 1024 //10M

var cache *freecache.Cache
var once sync.Once

func GetIns() *freecache.Cache {
	once.Do(func() {
		cache = freecache.NewCache(DefaultSize)
		debug.SetGCPercent(20)
	})

	return cache
}

func CacheSet(key string, value interface{}) error {
	return GetIns().Set([]byte(key), []byte(tool.ToJson(value)), 3600)
}

func CacheSetTtl(key string, value interface{}, ttl int) error {
	return GetIns().Set([]byte(key), []byte(tool.ToJson(value)), ttl)
}

func CacheGet(key string) (value []byte, err error) {
	return GetIns().Get([]byte(key))
}
