package factory

import (
	"fmt"
	"simple-http-server/store"
	"sync"
)

// StoreFactory 存储工厂
// 不同于字节暴露包级别的变量，封装成结构体更安全
type StoreFactory struct {
	sync.RWMutex                        // 读写锁
	providers    map[string]store.Store // 存储满足Store接口的实例
}

func NewStoreFactory() *StoreFactory {
	return &StoreFactory{providers: make(map[string]store.Store)}
}

func (f *StoreFactory) Register(name string, p store.Store) {
	f.Lock()
	defer f.Unlock()

	if p == nil {
		panic("Store: Register provider is nil")
	}
	// 检测是否重复注册
	if _, dup := f.providers[name]; dup {
		panic("Store: Register called twice for provider " + name)
	}
	// 注册
	f.providers[name] = p
}

// New 生产出/返回满足Store接口的实例
func (f *StoreFactory) New(providerName string) (store.Store, error) {
	f.RLock()
	p, ok := f.providers[providerName]
	f.RUnlock()
	if !ok {
		return nil, fmt.Errorf("store: unknown provider %s", providerName)
	}
	return p, nil
}
