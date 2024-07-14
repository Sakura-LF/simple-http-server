package store

import (
	"simple-http-server/store"
	"simple-http-server/store/factory"
	"sync"
)

// 图书数据存储模块的职责很清晰，就是用来存储整个 bookstore 的图书数据的。图书数据存
// 储有很多种实现方式，最简单的方式莫过于在内存中创建一个 map，以图书 id 作为 key，来
// 保存图书信息，

// 如果我们要考虑上生产环境，数据要进行持久化，那么最实际的方式就是通过 Nosql 数据库甚至是关系型数据库
// 实现对图书数据的存储与管理。

// 初始化模块
// 将 MemBookStore (图书数据存储模块)注册到工厂中
func init() {
	factory.Register("mem", &MemBookStore{Books: make(map[string]*store.Book)})
}

// 这里相当于数据库

// MemBookStore 是一个内存中的图书数据存储模块，它实现了 store.Store 接口。
type MemBookStore struct {
	sync.RWMutex
	Books map[string]*store.Book
}

// Create creates a new Book in the store.
func (ms *MemBookStore) Create(book *store.Book) error {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.Books[book.Id]; ok {
		return store.ErrExist
	}

	nBook := *book
	ms.Books[book.Id] = &nBook

	return nil
}

// Update updates the existed Book in the store.
func (ms *MemBookStore) Update(book *store.Book) error {
	ms.Lock()
	defer ms.Unlock()

	oldBook, ok := ms.Books[book.Id]
	if !ok {
		return store.ErrNotFound
	}

	nBook := *oldBook
	if book.Name != "" {
		nBook.Name = book.Name
	}

	if book.Authors != nil {
		nBook.Authors = book.Authors
	}

	if book.Press != "" {
		nBook.Press = book.Press
	}

	ms.Books[book.Id] = &nBook

	return nil
}

// Get retrieves a book from the store, by id. If no such id exists. an
// error is returned.
func (ms *MemBookStore) Get(id string) (store.Book, error) {
	ms.RLock()
	defer ms.RUnlock()

	t, ok := ms.Books[id]
	if ok {
		return *t, nil
	}
	return store.Book{}, store.ErrNotFound
}

// Delete deletes the book with the given id. If no such id exist. an error
// is returned.
func (ms *MemBookStore) Delete(id string) error {
	ms.Lock()
	defer ms.Unlock()

	// 如果没有这个 id，则返回 ErrNotFound 错误
	if _, ok := ms.Books[id]; !ok {
		return store.ErrNotFound
	}
	// 删除这个 id 对应的图书
	delete(ms.Books, id)
	return nil
}

// GetAll returns all the Books in the store, in arbitrary order.
func (ms *MemBookStore) GetAll() ([]store.Book, error) {
	ms.RLock()
	defer ms.RUnlock()

	// 创建一个切片，切片的长度等于 map 的长度，切片的容量等于 map 的长度
	allBooks := make([]store.Book, 0, len(ms.Books))
	for _, book := range ms.Books {
		allBooks = append(allBooks, *book)
	}
	return allBooks, nil
}
