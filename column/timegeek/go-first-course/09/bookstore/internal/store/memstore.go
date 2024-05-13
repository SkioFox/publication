package store

import (
	mystore "bookstore/store"
	factory "bookstore/store/factory"
	"sync"
)

func init() {
	factory.Register("mem", &MemStore{
		books: make(map[string]*mystore.Book),
	})
}

type MemStore struct {
	sync.RWMutex
	books map[string]*mystore.Book
}

// Create creates a new Book in the store.
func (ms *MemStore) Create(book *mystore.Book) error {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.books[book.Id]; ok {
		return mystore.ErrExist
	}
	/**
	如果ms.books[book.id] = book，那么由于存储的value是一个指针，当外部的book发生变化时，map中存储的value实际上也会跟着变化。
	我的代码中通过“clone”，实际上是将存储在map中的book与外部传入的book分离开来，它们是两块不同的存储区域。

	简而言之，这段代码首先通过解引用得到了book指针所指向的结构体的一个副本（赋值给了nBook），然后将这个副本的地址存入了一个映射中，作为特定键（book.Id）对应的值。
	这样的操作可能用在需要在不影响原结构体的同时，以某种方式修改或引用该结构体的一个副本的情景下，比如在实现深拷贝或者需要按ID管理结构体实例的场景。
	*/
	nBook := *book
	ms.books[book.Id] = &nBook

	return nil
}

// Update updates the existed Book in the store.
func (ms *MemStore) Update(book *mystore.Book) error {
	ms.Lock()
	defer ms.Unlock()

	oldBook, ok := ms.books[book.Id]
	if !ok {
		return mystore.ErrNotFound
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
func (ms *MemStore) Get(id string) (mystore.Book, error) {
	ms.RLock()
	defer ms.RUnlock()

	t, ok := ms.books[id]
	if ok {
		return *t, nil
	}
	return mystore.Book{}, mystore.ErrNotFound
}

// Delete deletes the book with the given id. If no such id exist. an error
// is returned.
func (ms *MemStore) Delete(id string) error {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.books[id]; !ok {
		return mystore.ErrNotFound
	}

	delete(ms.books, id)
	return nil
}

// GetAll returns all the books in the store, in arbitrary order.
func (ms *MemStore) GetAll() ([]mystore.Book, error) {
	ms.RLock()
	defer ms.RUnlock()
	// 注意，这里进行了深拷贝，确保返回的书籍切片中的每个书籍都是原始书籍的副本，外部对返回切片中的书籍所做的修改不会影响到MemStore内部的数据。
	allBooks := make([]mystore.Book, 0, len(ms.books))
	for _, book := range ms.books {

		allBooks = append(allBooks, *book)
	}
	return allBooks, nil
}
