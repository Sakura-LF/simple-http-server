package factory

import (
	"fmt"
	"simple-http-server/store"
	"sync"
)

//func NewFactory() *Factory {
//	return &Factory{
//		providers: make(map[string]store.Store),
//	}
//}

var (
	providersMu sync.RWMutex
	providers   = make(map[string]store.Store)
)

// Register 注册一个实现了Store接口的实例
func Register(name string, p store.Store) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if p == nil {
		panic("Store: Register provider is nil")
	}
	// 检测是否重复注册
	if _, dup := providers[name]; dup {
		panic("Store: Register called twice for provider " + name)
	}
	// 注册
	providers[name] = p
}

// New 生产出/返回满足Store接口的实例
func New(providerName string) (store.Store, error) {
	providersMu.RLock()
	p, ok := providers[providerName]
	providersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("store: unknown provider %s", providerName)
	}
	return p, nil
}
