package store

import (
	"simple-http-server/store"
	"simple-http-server/store/factory"
	"sync"
)

// 图书数据存储模块的职责很清晰，就是用来存储整个 bookstore 的图书数据的。图书数据存
// 储有很多种实现方式，最简单的方式莫过于在内存中创建一个 map，以图书 id 作为 key，来
// 保存图书信息，我们在这一讲中也会采用这种方式。但如果我们要考虑上生产环境，数据要进
// 行持久化，那么最实际的方式就是通过 Nosql 数据库甚至是关系型数据库，实现对图书数据
// 的存储与管理。

func init() {
	factory.Register("mem", &MemBookStore{books: make(map[string]*store.Book)})
}

// 这里相当于数据库

// MemBookStore 是一个内存中的图书数据存储模块，它实现了 store.Store 接口。
type MemBookStore struct {
	sync.RWMutex
	books map[string]*store.Book
}

// Create creates a new Book in the store.
func (ms *MemBookStore) Create(book *store.Book) error {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.books[book.Id]; ok {
		return store.ErrExist
	}

	nBook := *book
	ms.books[book.Id] = &nBook

	return nil
}

// Update updates the existed Book in the store.
func (ms *MemBookStore) Update(book *store.Book) error {
	ms.Lock()
	defer ms.Unlock()

	oldBook, ok := ms.books[book.Id]
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

	ms.books[book.Id] = &nBook

	return nil
}

// Get retrieves a book from the store, by id. If no such id exists. an
// error is returned.
func (ms *MemBookStore) Get(id string) (store.Book, error) {
	ms.RLock()
	defer ms.RUnlock()

	t, ok := ms.books[id]
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

	if _, ok := ms.books[id]; !ok {
		return store.ErrNotFound
	}

	delete(ms.books, id)
	return nil
}

// GetAll returns all the books in the store, in arbitrary order.
func (ms *MemBookStore) GetAll() ([]store.Book, error) {
	ms.RLock()
	defer ms.RUnlock()

	allBooks := make([]store.Book, 0, len(ms.books))
	for _, book := range ms.books {
		allBooks = append(allBooks, *book)
	}
	return allBooks, nil
}
