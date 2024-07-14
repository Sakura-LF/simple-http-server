package main

import (
	"simple-http-server/internal/store"
	store2 "simple-http-server/store"
	"simple-http-server/store/factory"
)

func main() {
	storeFactory := factory.NewStoreFactory()
	storeFactory.Register("memstore", &store.MemBookStore{Books: make(map[string]*store2.Book)})
}
