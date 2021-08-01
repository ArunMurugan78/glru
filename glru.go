package glru

import (
	"errors"
	"sync"

	"github.com/ArunMurugan78/glru/dll"
)

type Glru struct {
	nodeMap  map[string]*dll.Node
	maxItems int
	list     *dll.Dll
}

type Config struct {
	MaxItems int
}

var ErrKeyNotFound = errors.New("key not found")

func New(config Config) *Glru {
	return &Glru{maxItems: config.MaxItems, list: dll.New(), nodeMap: make(map[string]*dll.Node)}
}

func (cache *Glru) Set(key string, value interface{}) {
	var mutex sync.Mutex

	mutex.Lock()
	defer mutex.Unlock()

	node, ok := cache.nodeMap[key]

	if ok {
		if node.Value != value {
			node.Value = value
		}

		return
	}

	cache.nodeMap[key] = cache.list.Prepend(value)
}

func (cache *Glru) Get(key string) (interface{}, error) {
	var mutex sync.Mutex

	node, ok := cache.nodeMap[key]

	if ok {
		mutex.Lock()
		defer mutex.Unlock()

		// Brings the accessed node to the front of the list in O(1) time complexity
		cache.list.DeleteAndInsertAtHead(node)
		return node.Value, nil
	}
	return nil, ErrKeyNotFound
}
