package store

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrExist    = errors.New("exist")
)

type Book struct {
	Id      string   `json:"id"`      // 图书ISBN ID
	Name    string   `json:"name"`    // 图书名称
	Authors []string `json:"authors"` // 图书作者
	Press   string   `json:"press"`   // 出版社
}

type Store interface {
	Create(*Book) error
	Update(*Book) error
	/**
	这两个方法的返回值为什么不是*Book呢？查询出来的结果应该是Book的实例了，返回值类型是Book会发生Copy么？
	因为memstore底层使用的map类型为map[string]*mystore.Book， value用的是Book指针，如果直接返回指针有风险，这样外面的代码就可以直接修改memstore中的数据了。
	返回Book类型时会有一次mem copy。
	*/
	Get(string) (Book, error)
	GetAll() ([]Book, error)
	Delete(string) error
}
