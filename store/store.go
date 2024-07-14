package store

import "errors"

// 图书数据存储模块的职责很清晰，就是用来存储整个 bookstore 的图书数据的。图书数据存
// 储有很多种实现方式，最简单的方式莫过于在内存中创建一个 map，以图书 id 作为 key，来
// 保存图书信息，我们在这一讲中也会采用这种方式。但如果我们要考虑上生产环境，数据要进
// 行持久化，那么最实际的方式就是通过 Nosql 数据库甚至是关系型数据库，实现对图书数据
// 的存储与管理。

var (
	ErrNotFound = errors.New("not found")
	ErrExist    = errors.New("exist")
)

type Book struct {
	Id      string   `json:"id"`     // 图书ID
	Name    string   `json:"name"`   // 图书名称
	Authors []string `json:"author"` // 图书作者
	Press   string   `json:"press"`  // 出版社
}

// Store 商店接口
// 针对 Book 存取的接口类型 Store。
// 这样，对于想要进行图书数据操作的一方来说，他只需要得到一个满足 Store 接口的 实例，就可以实现对图书数据的存储操作了
// 不用再关心图书数据究竟采用了何种存储方式。
type Store interface {
	Create(*Book) error       // 创建新图书条目
	Update(*Book) error       // 更新某图书条目
	Get(string) (Book, error) // 获取单个图书
	GetAll() ([]Book, error)  // 获取全部图书
	Delete(string) error      // 删除某图书条目
}
