package tree

import (
	"fmt"
	"sync"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
	"github.com/pkg/errors"

	"github.com/ars0915/matching-system/entity"
)

var (
	ErrorPersonExist    = errors.New("person exist")
	ErrorPersonNotFound = errors.New("person not found")
)

type PersonTree struct {
	tree  *redblacktree.Tree
	idMap map[uint64]*entity.Person
	mu    sync.RWMutex
}

func NewPersonTree() *PersonTree {
	return &PersonTree{
		tree:  redblacktree.NewWith(utils.Float64Comparator),
		idMap: map[uint64]*entity.Person{},
	}
}

func (pt *PersonTree) AddPerson(p *entity.Person) error {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	_, exist := pt.idMap[p.ID]
	if exist {
		return ErrorPersonExist
	}

	pt.idMap[p.ID] = p
	if value, found := pt.tree.Get(p.Height); found {
		people := value.([]uint64)
		pt.tree.Put(p.Height, append(people, p.ID))
		return nil
	}
	pt.tree.Put(p.Height, []uint64{p.ID})

	return nil
}

func (pt *PersonTree) RemovePerson(id uint64) error {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	person, exist := pt.idMap[id]
	if !exist {
		return ErrorPersonNotFound
	}

	value, found := pt.tree.Get(person.Height)
	if !found {
		return ErrorPersonNotFound
	}

	delete(pt.idMap, id)

	ids := value.([]uint64)
	for i, v := range ids {
		if v != id {
			continue
		}
		ids = append(ids[:i], ids[i+1:]...)
		break
	}

	if len(ids) == 0 {
		pt.tree.Remove(person.Height)
		return nil
	}

	pt.tree.Put(person.Height, ids)
	return nil
}

func (pt *PersonTree) QueryByHeight(minHeight float64, maxHeight float64) []entity.Person {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	var result []entity.Person
	var ids []uint64

	iter := pt.tree.Iterator()
	for iter.Next() {
		height := iter.Key().(float64)
		if height >= minHeight && height <= maxHeight {
			ids = append(ids, iter.Value().([]uint64)...)
		}
	}

	for _, id := range ids {
		person, exist := pt.idMap[id]
		if exist {
			result = append(result, *person)
		}
	}

	return result
}

func (pt *PersonTree) FindByID(id uint64) (*entity.Person, bool) {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	person, exist := pt.idMap[id]
	return person, exist
}

func (pt *PersonTree) Print() {
	fmt.Println(pt.tree)
	fmt.Println(pt.tree.Values())
}
