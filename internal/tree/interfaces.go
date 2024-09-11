package tree

import "github.com/ars0915/matching-system/entity"

//go:generate mockgen -destination=../mocks/tree/person_tree.go -package=mocks github.com/ars0915/matching-system/internal/tree Tree
type (
	Tree interface {
		PersonTreeIface
	}
)

type (
	PersonTreeIface interface {
		AddPerson(p *entity.Person) error
		RemovePerson(id uint64) error
		QueryByHeight(minHeight float64, maxHeight float64) []entity.Person
		FindByID(id uint64) (*entity.Person, bool)
	}
)
