package main

import (
	"fmt"
	"context"
	"github.com/dgraph-io/ristretto"
)

var cache *ristretto.Cache

func main() {
	cache = initCache()
	jc := NewBridgeManager(context.Background(), db)
	jc.Listen()
}

func initCache() *ristretto.Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,     // number of keys to track frequency of (1M).
		MaxCost:     1 << 26, // maximum cost of cache (64MB).
		BufferItems: 64,      // number of keys per Get buffer.
		// Cost: func(value interface{}) int64 {
		// 	switch v := value.(type) {
		// 	case string:
		// 		return int64(len(v))
		// 	}
		// },
		OnEvict: func(item *ristretto.Item) {
			fmt.Println("Cache Eviction: ", item)
		},
	})

	if err != nil {
		panic(err)
	}

	return cache
}
